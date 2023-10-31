const timeWindow = 0.5 * 60 * 60 * 1000; // ms
const legend_cutoff = 40; // arc length pixels
const sincetext = " since last update"
const updateInterval = 1000; // ms

const units = "GW·s"
const unitscaling = 0.001

const month = new Intl.DateTimeFormat("en-US", {month: "short"})

// TODO: Set these dynamically based on axis text width/height
const x_offset = 50;
const y_offset = 30;

const data = {
    "regions": {}, // name: {name}
    "categories": {}, // name: {name, color}
    "latest": {}, // time, requirement, total, categories
    "periods": {
        "timestamps": [],
        "total": [],
        "requirement": [],
        "categories": [] // {name, inertia: []}
    }
};

function run(data) {

    const req = new XMLHttpRequest();
    req.open("GET", "/metadata", true);
    req.responseType = 'json';

    req.onload = function () {

        if (req.status != 200) {
            console.error(req.statusText);
            return;
        };

        data.regions = req.response.regions;
        data.categories = req.response.categories;

        ordered_cnames = Object.keys(data.categories)
            .sort((c1, c2) =>
                  data.categories[c1].order > data.categories[c2].order)

        data.periods.categories = ordered_cnames.map(c => {
            const category = data.categories[c]
            return { "name": category.name, inertia: [] };
        });

        update(data);
        setInterval(update, updateInterval, data);

    };

    req.onerror = function () {
        console.error(req.statusText);
    };

    req.send(null);

};

function update(data, last=0) {

    const req = new XMLHttpRequest();

    inertia_url = "/inertia?last="
    if (last) {
        inertia_url += last
    } else if (data.latest.time) {
        inertia_url += data.latest.time.valueOf()
    };

    req.open("GET", inertia_url, true);
    req.responseType = 'json';

    req.onload = function () {
        if (req.readyState === 4) {
            if (req.status === 200) {
                latest = req.response;
                data.lastpoll = new Date();
                data.latest = latest;
                appendInertia(data, latest);
                updateDisplay(data);
            } else if (req.status != 204) {
                console.error(req.statusText);
            };
        };
    };

    req.onerror = function () {
        console.error(req.statusText);
    };

    req.send(null);

};

function appendInertia(data, latest) {

    if (!Object.keys(data.periods.categories).length) {
        console.error("Appending to inertia before initialization");
    }

    data.periods.timestamps.push(latest.time);
    data.periods.requirement.push(latest.requirement);
    data.periods.total.push(latest.total.total_inertia);
    data.periods.categories.forEach(category => {
        category_inertia = latest.categories[category.name].total_inertia;
        category.inertia.push(category_inertia);
    });

};

function updateDisplay(data) {

    const latest = data.latest.time;

    updateText(data.latest);

    const frame = d3.select("#frame").node()
        .getBoundingClientRect();

    plotframe_w = frame.width * .67;
    ringframe_w = frame.width * .33;

    ring_x = plotframe_w + ringframe_w/2;
    ring_y = (frame.height - y_offset) / 2;
    ring_r = 0.7 * (Math.min(ringframe_w/2, ring_y) - y_offset);

    const legendroom = maxLegendWidth() + 40;

    var timeScale = d3.scaleTime()
        .domain([latest - timeWindow, latest])
        .range([x_offset, plotframe_w - legendroom]);

    d3.select("#t-axis")
        .call(d3.axisBottom(timeScale))
        .style("transform", "translate(0," + (frame.height - y_offset) + "px )");

    const inertia_max = maxInertia(data.periods);
    var inertiaScale = d3.scaleLinear()
        .domain([0, 1.1 * inertia_max])
        .range([frame.height - y_offset, 0]);

    d3.select("#i-axis")
        .call(d3.axisLeft(inertiaScale))
        .style("transform", "translate("+ x_offset + "px,0)");

    updateInertia(data, timeScale, inertiaScale);
    updateCategories(data, ring_x, ring_y, ring_r);
    updateRequirement(data, timeScale, inertiaScale);

    return latest

};

function updateCategories(data, x, y, r) {

    const categories = categoryRingData(data.latest, x, y, r);

    d3.select("#ring")
      .selectAll(".inertia-area")
      .data(categories)
      .join("path")
      .classed("inertia-area", true)
      .attr("stroke", d => data.categories[d.name].color)
      .attr("stroke-width", r * .7)
      .attr("d", d => d.path);

    d3.select("#ring")
      .selectAll(".inertia-legend")
      .data(categories)
      .join("text")
      .classed("inertia-legend", true)
      .style("display", d => (d.arclength > legend_cutoff) ? "" : "none" )
      .text(d => d.name)
      .attr("x", d => d.mid_x)
      .attr("y", d => d.mid_y - 10)
      .append("tspan")
      .text(d => (" (" + inertiaText(d.val) + ")"))
      .attr("x", d => d.mid_x)
      .attr("y", d => d.mid_y + 10);

}

function updateRequirement(data, timeScale, inertiaScale) {

    const ts = data.periods.timestamps.map(timeScale);
    const t_cutoff = timeScale.range()[0]

    d3.select("#requirement")
        .attr("points", makePoints(ts, data.periods.requirement.map(inertiaScale), t_cutoff));

}

function updateInertia(data, timeScale, inertiaScale) {

    const ts = data.periods.timestamps.map(timeScale);
    const t_cutoff = timeScale.range()[0]

    d3.select("#total-inertia")
        .attr("points", makePoints(ts, data.periods.total.map(inertiaScale), t_cutoff));

}

function updateText(currentData) {

    t = new Date(currentData.time);
    timestamp = t.getHours().toString().padStart(2, 0) + ":"
              + t.getMinutes().toString().padStart(2, 0) + ":"
              + t.getSeconds().toString().padStart(2, 0) + " "
              + month.format(t) + " " + t.getDate()

    d3.select("#time #lastupdate").text(timestamp);

    absolute = inertiaText(currentData.total.total_inertia);

    if (currentData.total.total_inertia > currentData.requirement) {

        surplus = currentData.total.total_inertia - currentData.requirement;
        relative = inertiaText(surplus) + " above threshold"
        color = "#EEEEEE";

    } else {
        shortfall = currentData.requirement - currentData.total.total_inertia;
        relative = inertiaText(shortfall) + " below threshold"
        color = "#FF0000";
    }

    d3.select("#inertia #absolute")
        .text(absolute)
        .style("color", color);

    d3.select("#inertia #relative")
        .text(relative)
        .style("color", color);

}

function makePoints(xs, ys, x_cutoff) {

    var result = ""

    for (var i = 0; i < xs.length; i++) {

        if (xs[i] < x_cutoff) { continue; }

        result += xs[i] + "," + ys[i] + " ";

    }

    return result

}

function makePolyPoints(xs, ylows, yhighs, x_cutoff) {

    var lows = ""
    var highs = ""

    for (var i = 0; i < xs.length; i++) {

        if (xs[i] < x_cutoff) { continue; }

        lows += xs[i] + "," + ylows[i] + " ";
        highs = xs[i] + "," + yhighs[i] + " " + highs;

    }

    return lows + highs

}

function maxInertia(periods) {

    var max = 0;

    for (const inertia of periods.total) {
        if (inertia > max) {
            max = inertia
        }
    }

    return max;

}

function inertiaText(x) {
    return (x * unitscaling).toFixed(1) + " " + units;
};

function categoryRingData(latest, cx, cy, r) {

    categories = [];
    cum_angle = 0

    for (const c of Object.keys(latest.categories)) {

        share = latest.categories[c].total_inertia / latest.total.total_inertia;
        angle = share * 2 * Math.PI;

        mid_x = cx + r * Math.sin(cum_angle + angle/2);
        mid_y = cy - r * Math.cos(cum_angle + angle/2);

        x = cx + r * Math.sin(cum_angle);
        y = cy - r * Math.cos(cum_angle);
        path = "M " + x + " " + y + " ";

        cum_angle += angle;

        x = cx + r * Math.sin(cum_angle);
        y = cy - r * Math.cos(cum_angle);

        path += "A " + r + " " + r + " 0 ";

        if (angle > Math.PI) {
            path += "1 1 "
        } else {
            path += "0 1 "
        }
        path += x + " " + y;

        categories.push({
            "name": c,
            "path": path,
            "mid_x": mid_x,
            "mid_y": mid_y,
            "val": latest.categories[c].total_inertia,
            "share": share,
            "arclength": r*angle
        });

    }

    return categories

}
function maxLegendWidth() {

    var legends = d3.select("#canvas")
      .selectAll(".inertia-legend")

    var max_width = 0;

    for (const node of legends.nodes()) {
      var w = node.getBoundingClientRect().width;
      if (w > max_width) {
          max_width = w
      }
    };

    return max_width;

};

function updateElapsed(data) {

    if (!data.lastpoll) {
        return
    }

    var seconds_elapsed = (Date.now() - data.lastpoll) / 1000

    const hours_elapsed = Math.floor(seconds_elapsed / 3600);
    seconds_elapsed %= 3600;

    const minutes_elapsed = Math.floor(seconds_elapsed / 60);
    seconds_elapsed %= 60;

    const elapsed = hours_elapsed.toString().padStart(2, 0) + ":"
        + minutes_elapsed.toString().padStart(2, 0) + ":"
        + Math.floor(seconds_elapsed).toString().padStart(2, 0)
        + sincetext;

    d3.select("#time #elapsed").text(elapsed);

}


// TODO: Trigger re-draw when window size changes
run(data);
setInterval(updateElapsed, 1000, data);
