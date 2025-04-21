import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt

filename = "olivia-pc"

# Load CSV data
df = pd.read_csv("data/" + filename + ".csv")

# Set up the subplot grid (2 rows, 2 columns)
fig, axs = plt.subplots(1, 2, figsize=(10, 5))
sns.set(style="whitegrid")

# --- Speedup across Array Sizes ---
# sns.lineplot(data=df, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs[1, 0])
# axs[1, 0].set_title("Speedup Across Array Sizes")

df_MM = df[df["Algorithm"].str.contains("MM")]
sns.lineplot(data=df_MM, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs[0])
axs[0].set_title("Matrix Multiplication Speedup")

# --- Speedup for MergeSort ---
df_MS = df[df["Algorithm"].str.contains("MS")]
sns.lineplot(data=df_MS, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs[1])
axs[1].set_title("Speedup for Merge Sort")

plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
plt.tight_layout()
plt.show()
