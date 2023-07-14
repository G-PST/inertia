# G-PST Inertia Monitoring Framework

![Test Status](https://github.com/G-PST/inertia/actions/workflows/tests.yml/badge.svg?branch=main)
[![Documentation](https://img.shields.io/badge/Documentation-latest-blue.svg)](https://g-pst.github.io/inertia)

`inertia` is a Go package for real-time estimation of
a power system's inertia levels. It defines software interfaces
for ingesting and reporting data in real-time.
Unit commitment ("H-constant")-based estimation logic and data interfaces
are available in the `inertia/uc` package. PMU-based estimation
methods are planned as future work.

System integrators can provide deployment-specfic data ingestion code
(e.g., developed for use with a specific EMS or historian system) that
conforms to the stated data interfaces for the desired estimation method.
Once these input interfaces are
implemented, ingested data can be automatically processed and reported out
via the package's real-time visualization framework.

This package provides two off-the-shelf visualization
modules in `inertia/viz/text` and `inertia/viz/web`, but
custom implementations of
the Visualizer interface can also be used. Multiple Visualizers can be
associated with a single real-time data stream, allowing for reporting to
multiple outputs at the same time, for example logging to a text file while
also visualizing results in a web browser.
