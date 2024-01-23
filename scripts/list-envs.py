#!/usr/bin/python3
# -*- compile-command: "./list-envs.sh"; -*-

import gymnasium as gym
for key in sorted(gym.envs.registry.keys()):
    print(key)
