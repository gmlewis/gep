import gymnasium as gym
# env = gym.make("LunarLander-v2", render_mode="human")
env = gym.make('Blackjack-v1', natural=False, sab=False, render_mode="human")

# Useful information:
print("Action space:", env.action_space)
print("Action space shape:", env.action_space.shape)
print("Observation space:", env.observation_space)
print("Observation space shape:", env.observation_space.shape)
from gymnasium.wrappers import FlattenObservation
wrapped_env = FlattenObservation(env)
print("Wrapped observation space:", wrapped_env.observation_space)
print("Wrapped observation space shape:", wrapped_env.observation_space.shape)

# Demonstrate totally random agent - no reinforcement learning:
observation, info = env.reset()
print("Initial observation:", observation)
print("Initial info:", info)

for _ in range(1000):
    action = env.action_space.sample()  # agent policy that uses the observation and info
    observation, reward, terminated, truncated, info = env.step(action)

    if terminated or truncated:
        observation, info = env.reset()

env.close()