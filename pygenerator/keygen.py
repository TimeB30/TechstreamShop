from ahk import AHK
import pyperclip
import sys
def set_days(num,ahk):
    win = ahk.win_get(title="Techstream Keygen v3.9")
    win.move(282, 65)
    win.activate()
    ahk.click(766,182,coord_mode="Screen")
    ahk.key_press("Backspace")
    ahk.key_press("Backspace")
    ahk.key_press("Backspace")
    ahk.key_press("Backspace")
    ahk.key_press("Backspace")
    ahk.type(str(num))
def set_version(index, ahk):
    win = ahk.win_get(title="Techstream Keygen v3.9")
    win.move(282, 65)
    win.activate()
    ahk.click(500, 182, coord_mode="Screen")
    for i in range(0, 50):
        ahk.key_press('Up')
    for i in range(0, index):
        ahk.key_press('Down')
    ahk.key_press('Enter')
def get_activator(ahk,text,version_index, days=365):
    return "test key"
def get_activator2(ahk,text,version_index, days=365):
    set_version(version_index, ahk)
    set_days(days,ahk)
    # ahk = AHK()
    count = 0
    win = ahk.win_get(title="Techstream Keygen v3.9")
    win.move(282, 65)
    win.activate()
    while(count < 5):
        count += 1
        answers = set()
        ahk.click(579, 132,coord_mode="Screen")
        ahk.click(579, 132,coord_mode="Screen")
        ahk.key_press("Backspace")
        ahk.type(text)
        ahk.click(755, 128,coord_mode="Screen")
        ahk.click(500, 248,coord_mode="Screen")
        ahk.click(500, 248,coord_mode="Screen")
        ahk.right_click()
        ahk.click(559, 312,coord_mode="Screen")
        result = "Japan: " + pyperclip.paste()
        answers.add(pyperclip.paste())
        ahk.click(505, 272,coord_mode="Screen")
        ahk.click(505, 272,coord_mode="Screen")
        ahk.right_click()
        ahk.click(549, 327,coord_mode="Screen")
        result += "\n\nNorth America: " + pyperclip.paste()
        answers.add(pyperclip.paste())
        ahk.click(484, 299,coord_mode="Screen")
        ahk.click(484, 299,coord_mode="Screen")
        ahk.right_click()
        ahk.click(520, 360,coord_mode="Screen")
        result += "\n\nEurope: " + pyperclip.paste()
        answers.add(pyperclip.paste())
        ahk.click(508, 327,coord_mode="Screen")
        ahk.click(508, 327,coord_mode="Screen")
        ahk.right_click()
        ahk.click(560, 386,coord_mode="Screen")
        result += "\n\nOther: " + pyperclip.paste()
        answers.add(pyperclip.paste())
        if len(answers) == 4:
            return result
    return ""

