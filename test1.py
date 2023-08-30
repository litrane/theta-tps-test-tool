def calculate_tx_num_sum(filename, start_height, end_height):
    total_tx_num = 0
    in_target_range = False

    with open(filename, 'r') as file:
        for line in file:
            if 'ProposeBlockTxs:' in line:
                if 'block.height =' in line:
                    parts = line.split('block.height =')
                    if len(parts) > 1:
                        height_str = parts[1].split()[0].strip()
                        try:
                            height = int(height_str)
                            if start_height <= height <= end_height:
                                in_target_range = True
                            else:
                                in_target_range = False
                        except ValueError:
                            pass

                if in_target_range and 'tx num =' in line:
                    parts = line.split('tx num =')
                    if len(parts) > 1:
                        tx_num_str = parts[1].strip()
                        try:
                            tx_num = int(tx_num_str)
                            total_tx_num += tx_num
                        except ValueError:
                            pass

    return total_tx_num


filename = 'output1'
start_height = 56
end_height = 110
sum_tx_num = calculate_tx_num_sum(filename, start_height, end_height)
print(
    f"Sum of tx num from block.height {start_height} to {end_height}: {sum_tx_num}")
