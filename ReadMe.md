# GoMandlebrodt

This project contains a go and a python file used to both generate and display various gifs of the mandlebrodt set respectively.

## Installation

Use the package manager [pip](https://pip.pypa.io/en/stable/) to install both [numpy](https://pypi.org/project/numpy/) and [array2gif](https://pypi.org/project/array2gif/).

```bash
pip install numpy
pip install array2gif
```
Also ensure that you are able to run powershell scripts on your machine. To do this you may have to run "set-executionpolicy unrestricted".

## Usage
Running the following command in powershell will create a gif called test.gif  in your current working directory.
```bash
.\BuildGif.ps1
```
Also note that a text file with information about each frame will be created in the current working directory to pass information between go and python.

To modify various parameters regarding what will eventually be in the final gif. Make sure that when you modify them you also have them consistent across both the go and the python files.

In the python file:
```python
#This is somewhere near the top
size_of_image = 1024
num_frames = 30
max_iterations = 1000

```
In the go file:
```go
#This is somewhere near the start of main
frame_dimension := float64(1024)
number_frames := float64(30)
max_iterations := 1000

```
**Note**: more than 100 total frames is not currently supported. It would not be difficult to change so that it does work however, but it may take a while for the program to complete.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)
