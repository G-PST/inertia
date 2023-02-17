var colors = {
    "Coal": "#222222",
    "NG-Steam": "#3D3376",
    "NG-CC": "#52216B",
    "NG-CT": "#C2A1DB",
    "Biomass": "#5B9844",
    "Petroleum": "#853D65"
};

var updates = [
    {
        "time": new Date("2022-01-01T00:00:00"), "requirement": 100,
        "data": { "Coal" : 30, "NG-CC": 40, "NG-CT": 40}
    }, {
        "time": new Date("2022-01-01T00:05:00"), "requirement": 100,
        "data": { "Coal" : 30, "NG-CC": 30, "NG-CT": 40}
    }, {
        "time": new Date("2022-01-01T00:10:00"), "requirement": 100,
        "inertia": { "Coal" : 30, "NG-CC": 30, "NG-CT": 30}
    }, {
        "time": new Date("2022-01-01T00:15:00"), "requirement": 100,
        "inertia": { "Coal" : 30, "NG-CC": 35, "NG-CT": 40}
    }
];

var timeWindow = 0.5 * 60 * 60 * 1000;

// regenerate when new data is received
const data = {
    "timestamps": [new Date("2022-01-01T00:00:00"),
                   new Date("2022-01-01T00:05:00"),
                   new Date("2022-01-01T00:10:00"),
                   new Date("2022-01-01T00:15:00")],
    "requirement": [100, 100, 100, 100],
    "inertia": [
        {"name": "Coal", "data": [30, 30, 30, 30]},
        {"name": "NG-CC", "data": [40, 30, 30, 35]},
        {"name": "NG-CT", "data": [40, 40, 30, 40]}]
};

var x_offset = 30;
var y_offset = 30;
const sincetext = " since last update"

function update(data) {

    const currentData = current(data);
    const latest = currentData.timestamp;

    updateText(currentData);

    const frame = d3.select("#frame").node()
        .getBoundingClientRect();

    var timeScale = d3.scaleTime()
        .domain([latest - timeWindow, latest])
        .range([x_offset,0.80 * frame.width]);


    d3.select("#t-axis")
        .call(d3.axisBottom(timeScale))
        .style("transform", "translate(0," + (frame.height - y_offset) + "px )");

    var inertia_max = maxInertia(data);
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

    const T = data.timestamps.length;
    const ts = data.timestamps.map(timeScale);

    var cum_inertia = new Array(T).fill(0)
    var cum_inertia_prev = new Array(T).fill(0)
    var categories = [];

    for (const category of data.inertia) {

        for (var t = 0; t < T; t++) {
            cum_inertia[t] = cum_inertia_prev[t] + category.data[t];
        }

        points = makePolyPoints(ts, 
            cum_inertia_prev.map(inertiaScale),
            cum_inertia.map(inertiaScale));

        categories.push({
            "name": category.name,
            "points": points,
            "mid": inertiaScale((cum_inertia[T-1] + cum_inertia_prev[T-1]) / 2),
            "val": category.data[T-1]
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
      .classed(d => d.name, true)
      .attr("fill", d => colors[d.name])
      .attr("points", d => d.points);

    d3.select("#canvas")
      .selectAll(".inertia-legend")
      .data(categories)
      .join("text")
      .classed("inertia-legend", true)
      .classed(d => d.name, true)
      .text(d => d.name + " (" + d.val + " GW·s)")
      .attr("x", "85%")
      .attr("y", d => d.mid);

}

function updateRequirement(data, timeScale, inertiaScale) {

    const ts = data.timestamps.map(timeScale);

    d3.select("#requirement")
        .attr("points", makePoints(ts, data.requirement.map(inertiaScale)));

}

function updateText(currentData) {

    t = currentData.timestamp;

    timestamp = t.getHours().toString().padStart(2, 0) + ":"
              + t.getMinutes().toString().padStart(2, 0) + ":"
              + t.getSeconds().toString().padStart(2, 0)

    d3.select("#time #lastupdate").text(timestamp);
    d3.select("#time #elapsed").text("00:00:00" + sincetext);

    absolute = currentData.inertia + " GW·s"

    if (currentData.inertia > currentData.requirement) {

        surplus = currentData.inertia - currentData.requirement;
        relative = surplus + " GW·s above threshold"
        color = "#EEEEEE";

    } else {
        shortfall = currentData.requirement - currentData.inertia;
        relative = shortfall + " GW·s below threshold"
        color = "#FF0000";
    }

    d3.select("#inertia #absolute")
        .text(absolute)
        .style("color", color);

    d3.select("#inertia #relative")
        .text(relative)
        .style("color", color);

}

function updateElapsed(latest) {

    var seconds_elapsed = (Date.now() - latest) / 1000

    const hours_elapsed = Math.floor(seconds_elapsed / 3600);
    seconds_elapsed %= 3600;

    const minutes_elapsed = Math.floor(seconds_elapsed / 60);
    seconds_elapsed %= 60;

    const elapsed = hours_elapsed.toString().padStart(2, 0) + ":"
        + minutes_elapsed.toString().padStart(2, 0) + ":" + Math.floor(seconds_elapsed).toString().padStart(2, 0)
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

function makePolyPoints(xs, ylows, yhighs) {

    var lows = ""
    var highs = ""

    for (var i = 0; i < xs.length; i++) {
        lows += xs[i] + "," + ylows[i] + " ";
        highs = xs[i] + "," + yhighs[i] + " " + highs;
    }

    return lows + highs

}

function current(data) {

    const last = data.timestamps.length - 1;
    const currentTime = data.timestamps[last];
    const currentRequirement = data.requirement[last];
    var currentInertia = 0;

    for (const category of data.inertia) {
        currentInertia += category.data[last];
    }

    return {
        "timestamp": currentTime,
        "requirement": currentRequirement,
        "inertia": currentInertia
    };
}

function maxInertia(data) {

    var max = 0;

    for (t = 0; t < data.timestamps.length; t++) {

        inertia = 0;

        for (const category of data.inertia) {
            inertia += category.data[t];
        }

        if (inertia > max) {
            max = inertia
        }

    }

    return max;

}

// Define history length (e.g. 24 hours): plot points
// from lastest - history to latest

// - Receive new inertia state from websocket
// - Pass to updater function, recalculate visualized set
//   based on history length
// - Regenerate time scale based on latest point and history length
// Parse times to collect min/max and convert to linear scale
// 

var latest = update(data);

setInterval(updateElapsed, 1000, latest);
