import re

# 定义正则表达式模式
pattern = r'.*map size.*'

# 打开输入文件和输出文件
with open('./outputmain.log', 'r') as input_file, open('out.txt', 'w') as output_file:
    # 逐行读取输入文件
    for line in input_file:
        # 使用正则表达式匹配每一行
        if re.match(pattern, line):
            # 如果匹配成功，将该行写入输出文件
            output_file.write(line)
