import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np
from scipy import stats

filename = "olivia-laptop-3"

# Load CSV data
df = pd.read_csv("data/" + filename + ".csv")

# Set up the subplot grid (2 rows, 2 columns)
fig, axs = plt.subplots(1, 1, figsize=(8, 5))
sns.set_style('whitegrid')

# --- Speedup across Array Sizes ---
# df_MM = df[df["Algorithm"].str.contains("MM")]
# sns.lineplot(data=df_MM, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs[0])
# axs[0].set_title("Matrix Multiplication Speedup")

# --- Speedup for MergeSort ---
df_MS = df[df["Algorithm"].str.contains("MS")]
sns.lineplot(data=df_MS, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs)
axs.set_title("Speedup for Merge Sort")

plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
plt.tight_layout()
plt.show()

base_algorithm = 'MS-Seq'

for array_size in df['ArraySizeA'].unique():
    print(f"\nProcessing ArraySizeA = {array_size}")

    size_group = df[df['ArraySizeA'] == array_size]
    
    if base_algorithm in size_group['Algorithm'].values:
        ms_seq_group = size_group[size_group['Algorithm'] == base_algorithm]
        print(f"  Found {base_algorithm} for ArraySizeA = {array_size}.")

        for algorithm in size_group['Algorithm'].unique():
            if algorithm != base_algorithm:
                other_group = size_group[size_group['Algorithm'] == algorithm]
                print(f"    Comparing {base_algorithm} vs {algorithm} for ArraySizeA = {array_size}")

                # Perform Shapiro-Wilk test for normality for MS-Seq and the other algorithm
                stat_ms_seq, p_value_ms_seq = stats.shapiro(ms_seq_group['AvgTime'])
                stat_other, p_value_other = stats.shapiro(other_group['AvgTime'])

                print(f"      Shapiro-Wilk for {base_algorithm}: statistic = {stat_ms_seq}, p-value = {p_value_ms_seq}")
                print(f"      Shapiro-Wilk for {algorithm}: statistic = {stat_other}, p-value = {p_value_other}")

                # Check normality and perform the appropriate test
                if p_value_ms_seq < 0.05 or p_value_other < 0.05:
                    # If either dataset is not normally distributed, use Mann-Whitney U test
                    print(f"        Data is NOT normally distributed. Using Mann-Whitney U test...")
                    u_stat, u_p_value = stats.mannwhitneyu(ms_seq_group['AvgTime'], other_group['AvgTime'], alternative='two-sided')
                    print(f"        Mann-Whitney U test result: U-statistic = {u_stat}, p-value = {u_p_value}")
                    if u_p_value < 0.05:
                        print("  ✅ Statistically significant difference (p < 0.05)")
                    else:
                        print("  ⚠️ No significant difference (p ≥ 0.05)")
                else:
                    # If both datasets are normal, use t-test
                    print(f"        Data is normally distributed. Using t-test...")
                    t_stat, t_p_value = stats.ttest_ind(ms_seq_group['AvgTime'], other_group['AvgTime'])
                    print(f"        t-test result: t-statistic = {t_stat}, p-value = {t_p_value}")
                    if t_p_value < 0.05:
                        print("  ✅ Statistically significant difference (p < 0.05)")
                    else:
                        print("  ⚠️ No significant difference (p ≥ 0.05)")


