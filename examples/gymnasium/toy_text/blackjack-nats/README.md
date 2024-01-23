# Blackjack Gymnasium

This example is taken from here:
https://gymnasium.farama.org/tutorials/training_agents/blackjack_tutorial/

## Installation

First, install python3 and anaconda3.
In a shell window, create a new conda environment and activate it:

```shell
conda create -n openai-gym python=3.10.8
conda activate openai-gym
pip install "gymnasium[all]"
pip install nats-py
```

## Benchmarking

To benchmark the all-Python example, here are the times from my Mac M2 Max:

```shell
$ time python3 blackjack_tutorial.py
100%|████████████████| 100000/100000 [00:05<00:00, 16992.24it/s]

real	0m6.733s
user	0m12.424s
sys	0m0.124s
```

To benchmark the Python + Go + NATS example, first start the Go program
in one terminal window:

```shell
go run main.go
```

Then, in a second terminal window, run the Python script:

```
$ time python3 blackjack-nats.py
100%|█████████████████| 100000/100000 [00:53<00:00, 1855.64it/s]

real	0m54.760s
user	0m22.380s
sys	0m8.624s
```