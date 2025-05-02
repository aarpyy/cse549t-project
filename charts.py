import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import numpy as np
from scipy import stats
import matplotlib.gridspec as gridspec

filename = "olivia-laptop-3"

# load data
df1 = pd.read_csv("data/" + filename + ".csv")
df2 = pd.read_csv("data/scott-laptop.csv")

# set up plots
fig = plt.figure(figsize=(14, 7))
gs = gridspec.GridSpec(2, 2, width_ratios=[1,1])

ax1 = fig.add_subplot(gs[0, 0])
ax2 = fig.add_subplot(gs[0, 1])
ax3 = fig.add_subplot(gs[1, 0])
ax4 = fig.add_subplot(gs[1,1])

sns.set_style('whitegrid')

# # --- Speedup across Array Sizes ---
# # df_MM = df[df["Algorithm"].str.contains("MM")]
# # sns.lineplot(data=df_MM, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs[0])
# # axs[0].set_title("Matrix Multiplication Speedup")

# # --- Speedup for MergeSort ---
# # df_MS = df[df["Algorithm"].str.contains("MS")]
# # sns.lineplot(data=df_MS, x="ArraySizeA", y="Speedup", hue="Algorithm", marker="o", ax=axs)
# # axs.set_title("Speedup for Merge Sort")


# get only same sizes
shared_sizes = set(df1["ArraySizeA"]).intersection(set(df2["ArraySizeA"]))
df1 = df1[df1["ArraySizeA"].isin(shared_sizes)]
df2 = df2[df2["ArraySizeA"].isin(shared_sizes)]

df1["Source"] = "Processor B"
df2["Source"] = "Processor C"

dfcat = pd.concat([df1, df2], ignore_index=True)

sns.lineplot(data=dfcat[dfcat["Algorithm"] == "MS-Seq"], x="ArraySizeA", y="AvgTime", hue="Source", ax=ax1, marker='o')
ax1.set_title("MS-Seq")
ax1.set_ylabel("Runtime (μs)")
ax1.set_xlabel("Array Size")

sns.lineplot(data=dfcat[dfcat["Algorithm"] == "MS-CacheAware"], x="ArraySizeA", y="AvgTime", hue="Source", ax=ax2, marker='o')
ax2.set_title("MS-CacheAware")
ax2.set_ylabel("Runtime (μs)")
ax2.set_xlabel("Array Size")

sns.lineplot(data=dfcat[dfcat["Algorithm"] == "MS-Par"], x="ArraySizeA", y="AvgTime", hue="Source", ax=ax3, marker='o')
ax3.set_title("MS-Par")
ax3.set_ylabel("Runtime (μs)")
ax3.set_xlabel("Array Size")

sns.lineplot(data=dfcat[dfcat["Algorithm"] == "MS-PMerge"], x="ArraySizeA", y="AvgTime", hue="Source", ax=ax4, marker='o')
ax4.set_title("MS-PMerge")
ax4.set_ylabel("Runtime (μs)")
ax4.set_xlabel("Array Size")

# base_algorithm = 'MS-Seq'

# for array_size in df['ArraySizeA'].unique():
#     print(f"\nProcessing ArraySizeA = {array_size}")

#     size_group = df[df['ArraySizeA'] == array_size]
    
#     if base_algorithm in size_group['Algorithm'].values:
#         ms_seq_group = size_group[size_group['Algorithm'] == base_algorithm]
#         print(f"  Found {base_algorithm} for ArraySizeA = {array_size}.")

#         for algorithm in size_group['Algorithm'].unique():
#             if algorithm != base_algorithm:
#                 other_group = size_group[size_group['Algorithm'] == algorithm]
#                 print(f"    Comparing {base_algorithm} vs {algorithm} for ArraySizeA = {array_size}")

#                 # Perform Shapiro-Wilk test for normality for MS-Seq and the other algorithm
#                 stat_ms_seq, p_value_ms_seq = stats.shapiro(ms_seq_group['AvgTime'])
#                 stat_other, p_value_other = stats.shapiro(other_group['AvgTime'])

#                 print(f"      Shapiro-Wilk for {base_algorithm}: statistic = {stat_ms_seq}, p-value = {p_value_ms_seq}")
#                 print(f"      Shapiro-Wilk for {algorithm}: statistic = {stat_other}, p-value = {p_value_other}")

#                 # Check normality and perform the appropriate test
#                 if p_value_ms_seq < 0.05 or p_value_other < 0.05:
#                     # If either dataset is not normally distributed, use Mann-Whitney U test
#                     print(f"        Data is NOT normally distributed. Using Mann-Whitney U test...")
#                     u_stat, u_p_value = stats.mannwhitneyu(ms_seq_group['AvgTime'], other_group['AvgTime'], alternative='two-sided')
#                     print(f"        Mann-Whitney U test result: U-statistic = {u_stat}, p-value = {u_p_value}")
#                     if u_p_value < 0.05:
#                         print("  ✅ Statistically significant difference (p < 0.05)")
#                     else:
#                         print("  ⚠️ No significant difference (p ≥ 0.05)")
#                 else:
#                     # If both datasets are normal, use t-test
#                     print(f"        Data is normally distributed. Using t-test...")
#                     t_stat, t_p_value = stats.ttest_ind(ms_seq_group['AvgTime'], other_group['AvgTime'])
#                     print(f"        t-test result: t-statistic = {t_stat}, p-value = {t_p_value}")
#                     if t_p_value < 0.05:
#                         print("  ✅ Statistically significant difference (p < 0.05)")
#                     else:
#                         print("  ⚠️ No significant difference (p ≥ 0.05)")

column_name = 'AvgTime'
array_size_column = 'ArraySizeA'  # Update this with the correct column name if different
algorithm_column = 'Algorithm'   # Update this with the correct column name if different

# # Loop over unique combinations of ArraySizeA and Algorithm
# for array_size in df1[array_size_column].unique():
#     print(f"\nProcessing ArraySizeA = {array_size}")
    
#     # Filter data for the current ArraySizeA
#     df1_size_group = df1[df1[array_size_column] == array_size]
#     df2_size_group = df2[df2[array_size_column] == array_size]
    
#     # For each algorithm, compare runtimes between df1 and df2
#     for algorithm in df1_size_group[algorithm_column].unique():
#         print(f"  Comparing Algorithm = {algorithm} for ArraySizeA = {array_size}")
        
#         # Filter data for the current algorithm
#         df1_algo_group = df1_size_group[df1_size_group[algorithm_column] == algorithm]
#         df2_algo_group = df2_size_group[df2_size_group[algorithm_column] == algorithm]

#         # Perform Shapiro-Wilk test for normality for both df1 and df2
#         stat_df1, p_value_df1 = stats.shapiro(df1_algo_group[column_name])
#         stat_df2, p_value_df2 = stats.shapiro(df2_algo_group[column_name])

#         # Output the Shapiro-Wilk test results
#         print(f"    Shapiro-Wilk for df1 (Algorithm = {algorithm}): statistic = {stat_df1}, p-value = {p_value_df1}")
#         print(f"    Shapiro-Wilk for df2 (Algorithm = {algorithm}): statistic = {stat_df2}, p-value = {p_value_df2}")

#         # Check if either dataset is not normally distributed and select the appropriate test
#         if p_value_df1 < 0.05 or p_value_df2 < 0.05:
#             # If either dataset is not normally distributed, use Mann-Whitney U test
#             print(f"    Data is NOT normally distributed. Using Mann-Whitney U test...")
#             u_stat, u_p_value = stats.mannwhitneyu(df1_algo_group[column_name], df2_algo_group[column_name], alternative='two-sided')
#             print(f"    Mann-Whitney U test result: U-statistic = {u_stat}, p-value = {u_p_value}")
#             if u_p_value < 0.05:
#                 print(f"    ✅ Statistically significant difference (p < 0.05)")
#             else:
#                 print(f"    ⚠️ No significant difference (p ≥ 0.05)")
#         else:
#             # If both datasets are normally distributed, use t-test
#             print(f"    Data is normally distributed. Using t-test...")
#             t_stat, t_p_value = stats.ttest_ind(df1_algo_group[column_name], df2_algo_group[column_name])
#             print(f"    t-test result: t-statistic = {t_stat}, p-value = {t_p_value}")
#             if t_p_value < 0.05:
#                 print(f"    ✅ Statistically significant difference (p < 0.05)")
#                 ax4.annotate(
#                 f'Significant\n(p < 0.05)', 
#                 xy=(array_size, size_group['AvgTime'].mean()),  # Position to place the annotation
#                 xytext=(array_size, size_group['AvgTime'].mean() + 500),  # Adjust the offset
#                 arrowprops=dict(arrowstyle="->", color='red'),
#                 color='red',
#                 fontsize=10,
#                 ha='center'
#             )
#             else:
#                 print(f"    ⚠️ No significant difference (p ≥ 0.05)")

# # Tidy up
# plt.tight_layout()
# plt.show()

# Loop over unique combinations of ArraySizeA and Algorithm
for array_size in df1[array_size_column].unique():
    print(f"\nProcessing ArraySizeA = {array_size}")
    
    # Filter data for the current ArraySizeA
    df1_size_group = df1[df1[array_size_column] == array_size]
    df2_size_group = df2[df2[array_size_column] == array_size]

    # For each algorithm, compare runtimes between df1 and df2
    for algorithm in df1_size_group[algorithm_column].unique():
        print(f"  Comparing Algorithm = {algorithm} for ArraySizeA = {array_size}")
        
        # Filter data for the current algorithm
        df1_algo_group = df1_size_group[df1_size_group[algorithm_column] == algorithm]
        df2_algo_group = df2_size_group[df2_size_group[algorithm_column] == algorithm]

        stat_df1, p_value_df1 = stats.shapiro(df1_algo_group[column_name])
        stat_df2, p_value_df2 = stats.shapiro(df2_algo_group[column_name])

        print(f"    Shapiro-Wilk for df1 (Algorithm = {algorithm}): statistic = {stat_df1}, p-value = {p_value_df1}")
        print(f"    Shapiro-Wilk for df2 (Algorithm = {algorithm}): statistic = {stat_df2}, p-value = {p_value_df2}")

        # Check if either dataset is not normally distributed and select the appropriate test
        if p_value_df1 < 0.05 or p_value_df2 < 0.05:
            # If either dataset is not normally distributed, use Mann-Whitney U test
            print(f"    Data is NOT normally distributed. Using Mann-Whitney U test...")
            u_stat, u_p_value = stats.mannwhitneyu(df1_algo_group[column_name], df2_algo_group[column_name], alternative='two-sided')
            print(f"    Mann-Whitney U test result: U-statistic = {u_stat}, p-value = {u_p_value}")
            if u_p_value < 0.05:
                print(f"    ✅ Statistically significant difference (p < 0.05)")

                # Annotate the plot with the statistical result
                ax = {'MS-Seq': ax1, 'MS-CacheAware': ax2, 'MS-Par': ax3, 'MS-PMerge': ax4}.get(algorithm)
                y_min, y_max = ax.get_ylim()
                offset = (y_max - y_min) * 0.15
                ax.annotate(
                    f'p='+np.format_float_scientific(u_p_value, precision=4),
                    xy=(array_size, df1_algo_group[column_name].mean()),
                    xytext=(array_size, df1_algo_group[column_name].mean() + offset),
                    arrowprops=dict(arrowstyle="->", color='green'),
                    color='green',
                    fontsize=10,
                    ha='center'
                )

            else:
                print(f"    ⚠️ No significant difference (p ≥ 0.05)")
                ax = {'MS-Seq': ax1, 'MS-CacheAware': ax2, 'MS-Par': ax3, 'MS-PMerge': ax4}.get(algorithm)
                y_min, y_max = ax.get_ylim()
                offset = (y_max - y_min) * 0.15
                ax.annotate(
                    f'p='+np.format_float_scientific(u_p_value, precision=4),
                    xy=(array_size, df1_algo_group[column_name].mean()),  # Position to place the annotation
                    xytext=(array_size, df1_algo_group[column_name].mean() + offset),
                    arrowprops=dict(arrowstyle="->", color='red'),
                    color='red',
                    fontsize=10,
                    ha='center'
                )
        else:
            # If both datasets are normally distributed, use t-test
            print(f"    Data is normally distributed. Using t-test...")
            t_stat, t_p_value = stats.ttest_ind(df1_algo_group[column_name], df2_algo_group[column_name])
            print(f"    t-test result: t-statistic = {t_stat}, p-value = {t_p_value}")
            if t_p_value < 0.05:
                print(f"    ✅ Statistically significant difference (p < 0.05)")

                # Annotate the plot with the statistical result
                ax = {'MS-Seq': ax1, 'MS-CacheAware': ax2, 'MS-Par': ax3, 'MS-PMerge': ax4}.get(algorithm)
                y_min, y_max = ax.get_ylim()
                offset = (y_max - y_min) * 0.15
                ax.annotate(
                    f'p='+np.format_float_scientific(t_p_value, precision=4),
                    xy=(array_size, df1_algo_group[column_name].mean()),
                    xytext=(array_size, df1_algo_group[column_name].mean() + offset),
                    arrowprops=dict(arrowstyle="->", color='green'),
                    color='green',
                    fontsize=10,
                    ha='center'
                )

            else:
                print(f"    ⚠️ No significant difference (p ≥ 0.05)")
                ax = {'MS-Seq': ax1, 'MS-CacheAware': ax2, 'MS-Par': ax3, 'MS-PMerge': ax4}.get(algorithm)
                y_min, y_max = ax.get_ylim()
                offset = (y_max - y_min) * 0.15
                ax.annotate(
                    f'p='+np.format_float_scientific(t_p_value, precision=4),
                    xy=(array_size, df1_algo_group[column_name].mean()),
                    xytext=(array_size, df1_algo_group[column_name].mean() + offset),
                    arrowprops=dict(arrowstyle="->", color='red'),
                    color='red',
                    fontsize=10,
                    ha='center'
                )

# Tidy up
plt.tight_layout()
plt.show()
