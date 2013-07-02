# gohistogram - Histograms in Go

![build status](https://circleci.com/gh/VividCortex/gohistogram.png?circle-token=d37ec652ea117165cd1b342400a801438f575209)

The histograms in this package are based on the algorithms found in
Ben-Haim & Tom-Tov's *A Streaming Parallel Decision Tree Algorithm*
([PDF](http://jmlr.org/papers/volume11/ben-haim10a/ben-haim10a.pdf)).
Histogram bins do not have a preset size. As values stream into
the histogram, bins are dynamically added and merged.

Another implementation can be found in the Apache Hive project (see
[NumericHistogram](http://hive.apache.org/docs/r0.11.0/api/org/apache/hadoop/hive/ql/udf/generic/NumericHistogram.html)).

An example:

![histogram](http://i.imgur.com/5OplaRs.png)

The accurate method of calculating quantiles (like percentiles) requires
data to be sorted. Streaming histograms make it possible to approximate
quantiles without sorting (or even individually storing) values.

NumericHistogram is the more basic implementation of a streaming
histogram. WeightedHistogram implements bin values as exponentially-weighted
moving averages.

A maximum bin size is passed as an argument to the constructor methods. A
larger bin size yields more accurate approximations at the cost of increased
memory utilization and performance.

A picture of kittens:

![stack of kittens](http://i.imgur.com/QxRTWAE.jpg)

## Getting started

Get the code into your workspace:

    $ cd $GOPATH
    $ git clone git@github.com:VividCortex/gohistogram.git ./src/github.com/VividCortex/gohistogram

You can run the tests now:

    $ cd src/github.com/VividCortex/gohistogram
    $ go test .

## API Documentation

Full source documentation can be found [here][godoc].

[godoc]: http://godoc.org/github.com/VividCortex/gohistogram

## License

Copyright (c) 2013 VividCortex

Released under MIT License. Check `LICENSE` file for details.
