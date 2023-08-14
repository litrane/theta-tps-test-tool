def calculate_sum_from_log(filename):
    total_sum = 0

    with open(filename, 'r') as file:
        for line in file:
            if 'maxNumsTxs:' in line:
                parts = line.split('maxNumsTxs:')
                if len(parts) > 1:
                    num_str = parts[1].strip()
                    try:
                        num = int(num_str)
                        total_sum += num
                    except ValueError:
                        pass

    return total_sum


filename = 'output.log'
sum_maxNumsTxs = calculate_sum_from_log(filename)
print(f"Sum of maxNumsTxs: sum_maxNumsTxs={sum_maxNumsTxs}")
