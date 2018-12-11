import numpy as np
from array2gif import write_gif

# See documentation on modifying the following. Make sure that these are consistent across the go file too
size_of_image = 1024
num_frames = 30
max_iterations = 1000

# This should not be changed as it is used in keeping track of which color to print for various points
threshold = max_iterations // 2

# This function converts a number of iterations to an RGB scale. Modify it if you want a differing color scheme.
def iterations_to_RGB(iterations):
    #if the point converges it has -1 iterations
    if iterations == -1:
        return [0, 0, 0]
	#if the point did not get calculated, it has 0 iterations
    if iterations == 0:
        return [255, 60, 200]
	#otherwise scale from black to red if it is far away from the edge of the set
    if iterations < threshold:
        return [(iterations // 4)*(threshold // 127), 0, 0]
	#scale from red to white if it is closer to the set
    else:
        return[255, ((iterations - threshold)//4)*(threshold // 127), ((iterations - threshold)//4)*(threshold // 127)]

# This function reads all the information from a certain file containing information about a frame and appends it to the current gif
def read_data(file_name, gif):
    input_data = {} #This dictionary will store all of the current points with the tuple (x,y) being the key and iterations being the value
    output_data = np.zeros((size_of_image, size_of_image, 3), dtype=np.uint8) #this will eventually be appended to gif as the current "frame"
    with open(file_name) as f:
        for line in f:
            x = line.split(", ")
            input_data[(float(x[0]), float(x[1]))] = float(x[2]) #assign each x,y coordinate to have it's iteration value in the dictionary
        # The top left corner of the gif is the [0,0] element, and so by iterating over the data in this fashion we draw the image correctly
        cur_x, cur_y = 0, size_of_image - 1
        for key in sorted(input_data.keys()): #go through all the x,y values in "order", that being sorted by x in ascending order, then sorted by y in ascending order
            output_data[cur_x, cur_y] = iterations_to_RGB(input_data[key]) #assign the RGB value to the right pixel
			#iteration math...
            cur_y -= 1
            if cur_y == -1:
                cur_y = size_of_image - 1
                cur_x += 1
    return output_data

#the final gif
gif = []

#for every frame data file:
for i in range(0, num_frames, 1):
    file_name = "frame" + str(i).zfill(2) + ".txt"
    gif.append(read_data(file_name, gif)) #append the frame to the gif

write_gif(gif, 'test.gif', fps=10) #save the gif!

