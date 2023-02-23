const timeWindow = 24 * 60 * 60 * 1000; // ms
const legend_cutoff = 8; // px
const sincetext = " since last update"
const updateInterval = 1000; // ms

const units = "GW·s"
const unitscaling = 0.001

// TODO: Category sort data
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
        data.periods.categories = Object.keys(data.categories).map(c => {
            const category = data.categories[c]
            return { "name": category.name, inertia: [] };
        });

        setInterval(update, updateInterval, data);

    };

    req.onerror = function () {
        console.error(req.statusText);
    };

    req.send(null);

};

function update(data, last=0) {

    updateElapsed(data.latest.time);

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
    data.periods.total.push(latest.total);
    data.periods.categories.forEach(category => {
        category_inertia = latest.inertia[category.name];
        category.inertia.push(category_inertia);
    });

};

function updateDisplay(data) {

    const latest = data.latest.time;

    updateText(data.latest);

    const frame = d3.select("#frame").node()
        .getBoundingClientRect();

    const legendroom = maxLegendWidth() + 40;

    var timeScale = d3.scaleTime()
        .domain([latest - timeWindow, latest])
        .range([x_offset, frame.width - legendroom]);


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

    updateCategories(data, timeScale, inertiaScale);
    updateRequirement(data, timeScale, inertiaScale);

    return latest

};

function updateCategories(data, timeScale, inertiaScale) {

    const plot_max_x = timeScale.range()[1];
    const categories = categoryPlotData(
        data.periods, timeScale, inertiaScale);

    d3.select("#canvas")
      .selectAll(".inertia-area")
      .data(categories)
      .join("polygon")
      .classed("inertia-area", true)
      .attr("fill", d => data.categories[d.name].color)
      .attr("points", d => d.points);

    d3.select("#canvas")
      .selectAll(".inertia-legend-swatch")
      .data(categories)
      .join("polygon")
      .classed("inertia-legend-swatch", true)
      .attr("points", "0,0 10,5 10,-5")
      .style("display", d => (d.height > legend_cutoff) ? "" : "none" )
      .attr("transform", d => ("translate(" + (plot_max_x + 10) + " " + d.mid + ")"))
      .style("fill", d => data.categories[d.name].color);

    d3.select("#canvas")
      .selectAll(".inertia-legend")
      .data(categories)
      .join("text")
      .classed("inertia-legend", true)
      .style("display", d => (d.height > legend_cutoff) ? "" : "none" )
      .text(d => d.name)
      .attr("x", plot_max_x + 30)
      .attr("y", d => d.mid + 5)
      .append("tspan")
      .text(d => (" (" + inertiaText(d.val) + ")"));

}

function updateRequirement(data, timeScale, inertiaScale) {

    const ts = data.periods.timestamps.map(timeScale);
    const t_cutoff = timeScale.range()[0]

    d3.select("#requirement")
        .attr("points", makePoints(ts, data.periods.requirement.map(inertiaScale), t_cutoff));

}

function updateText(currentData) {

    const t = new Date(currentData.time);

    timestamp = t.getHours().toString().padStart(2, 0) + ":"
              + t.getMinutes().toString().padStart(2, 0) + ":"
              + t.getSeconds().toString().padStart(2, 0)

    d3.select("#time #lastupdate").text(timestamp);

    absolute = inertiaText(currentData.total);

    if (currentData.total > currentData.requirement) {

        surplus = currentData.total - currentData.requirement;
        relative = inertiaText(surplus) + " above threshold"
        color = "#EEEEEE";

    } else {
        shortfall = currentData.requirement - currentData.total;
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

function updateElapsed(latest_time) {

    if (!latest_time) {
        return
    }

    var seconds_elapsed = (Date.now() - latest_time) / 1000

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

function categoryPlotData(periods, timeScale, inertiaScale) {

    const T = data.periods.timestamps.length;
    const ts = data.periods.timestamps.map(timeScale);
    const t_cutoff = timeScale.range()[0];

    var cum_inertia = new Array(T).fill(0);
    var cum_inertia_prev = new Array(T).fill(0);
    var categories = [];

    for (const category of periods.categories) {

        for (var t = 0; t < T; t++) {
            cum_inertia[t] = cum_inertia_prev[t] + category.inertia[t];
        }

        points = makePolyPoints(ts, 
            cum_inertia_prev.map(inertiaScale),
            cum_inertia.map(inertiaScale),
            t_cutoff
        );

        category_plotdata = {
            "name": category.name,
            "points": points,
            "height": inertiaScale(cum_inertia_prev[T-1]) - inertiaScale(cum_inertia[T-1]),
            "mid": inertiaScale((cum_inertia[T-1] + cum_inertia_prev[T-1]) / 2),
            "val": category.inertia[T-1]
        };

        categories.push(category_plotdata);

        for (var t = 0; t < T; t++) {
            cum_inertia_prev[t] = cum_inertia[t];
        }

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

// TODO: Trigger re-draw when window size changes
run(data);
