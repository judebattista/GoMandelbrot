import numpy as np
from array2gif import write_gif

size_of_image = 2048
num_frames = 30
max_iterations = 1000
threshold = max_iterations / 2
def iterations_to_RGB(iterations):
    #if the point converges it has -1 iterations
    if iterations == -1:
    #if the point did not get calculated, it has 0 iterations
        return [0, 0, 0]
    if iterations == 0:
        return [255, 60, 200]
    """
    for z in checks:
        if iterations < z:
            index = checks.index(z) + 1
            return [index * 255/len(checks), , ]
    return [255, 0, 0]
    """
    if iterations < threshold:
        return [iterations*(threshold / 255), 0, 0]
    else:
        return[255, (iterations - threshold)*(threshold / 255), (iterations - threshold)*(threshold / 255)]

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
            output_data[cur_x, cur_y] = iterations_to_RGB(input_data[key])
            cur_y -= 1
            if cur_y == -1:
                cur_y = size_of_image - 1
                cur_x += 1
    return output_data

gif = []

#for i in range(num_frames - 1, -1, -1):
for i in range(0, num_frames, 1):
    file_name = "frame" + str(i).zfill(2) + ".txt"
    gif.append(read_data(file_name, gif))

write_gif(gif, 'test.gif', fps=1)

