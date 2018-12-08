import numpy as np
from array2gif import write_gif

size_of_image = 512
num_frames = 30

def iterations_to_RGB(iterations):
    if iterations == 0:
        return [0, 0, 0]
    checks = [pow(2, x) for x in range(0, 10)]
    for z in checks:
        if iterations < z:
            index = checks.index(z) + 1
            return [index * 255/len(checks), 150, 150]
    return [255, 255, 255]



def read_data(file_name, gif):
    input_data = {}
    output_data = np.zeros((size_of_image, size_of_image, 3), dtype=np.uint8)
    with open(file_name) as f:
        for line in f:
            x = line.split(", ")
            input_data[(float(x[0]), float(x[1]))] = float(x[2])
        # The top left corner of the gif is the [0,0] element
        cur_x, cur_y = 0, size_of_image - 1
        for key in sorted(input_data.keys()):
           output_data[cur_x,cur_y] = iterations_to_RGB(input_data[key])
           cur_x += 1
           if cur_x == size_of_image:
               cur_x = 0
               cur_y -= 1
    return output_data

gif = []

for i in range(num_frames - 1, -1, -1):
    file_name = "frame" + str(i).zfill(2) + ".txt"
    gif.append(read_data(file_name, gif))

write_gif(gif, 'test.gif', fps=5)

