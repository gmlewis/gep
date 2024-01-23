# Blackjack Gymnasium

This example is taken from here:
https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/

The difference between this example and the one in the sibbling
`blackjack-nats` directory is that this one uses a socket to
communicate with Python instead of NATS.

## Installation

First, install python3 and anaconda3.
In a shell window, create a new conda environment and activate it:

```shell
conda create -n openai-gym python=3.10.8
conda activate openai-gym
pip install "gymnasium[all]"
```

## Benchmarking

To benchmark the example (using Python + socket + Go), first run
the Python script from the repo: github.com/gmlewis/gym-socket-api
in one terminal window:

```shell
python .
```

Then, in a second terminal window, run the Go program

```shell
time go run main.go
...
Processed 1002 steps in 6.411958375 seconds (156.27050916405864 steps/sec)

real	0m7.217s
user	0m0.248s
sys	0m0.301s
```

Performance is still obviously a major concern when communicating with
Gymnasium.
