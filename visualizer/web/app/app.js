const timeWindow = 0.5 * 60 * 60 * 1000;
const legend_min_inertia = 10;
const sincetext = " since last update"

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

function initialize(data) {
    initialize_metadata(data);
    update(data, "");
};

function initialize_metadata(data) {

    const req = new XMLHttpRequest();
    req.open("GET", "/metadata", true);
    req.responseType = 'json';

    req.onload = function () {
        if (req.readyState === 4) {
            if (req.status === 200) {
                data.regions = req.response.regions;
                data.categories = req.response.categories;
                data.periods.categories = Object.keys(data.categories).map(c => {
                    const category = data.categories[c]
                    return { "name": category.name, inertia: [] };
                });
            } else {
                console.error(req.statusText);
            }
        }
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

    data.periods.timestamps.push(latest.time);
    data.periods.requirement.push(latest.requirement);
    data.periods.total.push(latest.total);
    data.periods.categories.forEach(category => {
        category_inertia = latest.inertia[category.name]
        category.inertia.push(category_inertia);
    });

};

function updateDisplay(data) {

    const latest = data.latest.time;

    updateText(data.latest);

    const frame = d3.select("#frame").node()
        .getBoundingClientRect();

    var timeScale = d3.scaleTime()
        .domain([latest - timeWindow, latest])
        .range([x_offset, 0.80 * frame.width]);


    d3.select("#t-axis")
        .call(d3.axisBottom(timeScale))
        .style("transform", "translate(0," + (frame.height - y_offset) + "px )");

    var inertia_max = maxInertia(data.periods);
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

    const T = data.periods.timestamps.length;
    const ts = data.periods.timestamps.map(timeScale);

    var cum_inertia = new Array(T).fill(0)
    var cum_inertia_prev = new Array(T).fill(0)
    var categories = [];

    for (const category of data.periods.categories) {

        for (var t = 0; t < T; t++) {
            cum_inertia[t] = cum_inertia_prev[t] + category.inertia[t];
        }

        points = makePolyPoints(ts, 
            cum_inertia_prev.map(inertiaScale),
            cum_inertia.map(inertiaScale));

        categories.push({
            "name": category.name,
            "points": points,
            "mid": inertiaScale((cum_inertia[T-1] + cum_inertia_prev[T-1]) / 2),
            "val": category.inertia[T-1]
        });

        for (var t = 0; t < T; t++) {
            cum_inertia_prev[t] = cum_inertia[t];
        }

    }

    d3.select("#canvas")
      .selectAll(".inertia-area")
      .data(categories)
      .join("polygon")
      .classed("inertia-area", true)
      .attr("fill", d => data.categories[d.name].color)
      .attr("points", d => d.points);

    d3.select("#canvas")
      .selectAll(".inertia-legend")
      .data(categories)
      .join("text")
      .classed("inertia-legend", true)
      .style("display", d => (d.val > legend_min_inertia) ? "" : "none" )
      .text(d => d.name + " (" + d.val + " MW·s)")
      .attr("x", "85%")
      .attr("y", d => d.mid);

}

function updateRequirement(data, timeScale, inertiaScale) {

    const ts = data.periods.timestamps.map(timeScale);

    d3.select("#requirement")
        .attr("points", makePoints(ts, data.periods.requirement.map(inertiaScale)));

}

function updateText(currentData) {

    const t = new Date(currentData.time);

    timestamp = t.getHours().toString().padStart(2, 0) + ":"
              + t.getMinutes().toString().padStart(2, 0) + ":"
              + t.getSeconds().toString().padStart(2, 0)

    d3.select("#time #lastupdate").text(timestamp);

    absolute = currentData.total + " MW·s"

    if (currentData.total > currentData.requirement) {

        surplus = currentData.total - currentData.requirement;
        relative = surplus + " MW·s above threshold"
        color = "#EEEEEE";

    } else {
        shortfall = currentData.requirement - currentData.total;
        relative = shortfall + " MW·s below threshold"
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

function makePoints(xs, ys) {

    var result = ""

    for (var i = 0; i < xs.length; i++) {
        result += xs[i] + "," + ys[i] + " ";
    }

    return result

}

// TODO: Don't draw values outside the cutoff window
function makePolyPoints(xs, ylows, yhighs) {

    var lows = ""
    var highs = ""

    for (var i = 0; i < xs.length; i++) {
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

// TODO: Trigger re-draw when window size changes
initialize(data);
setInterval(update, 1000, data);
