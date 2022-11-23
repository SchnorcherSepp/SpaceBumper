#!/usr/bin/env python

"""
This file contains an example in Python for an AI controlled client.
Use this example to program your own AI in Python.
"""

import socket
import time
from threading import Thread
from threading import Lock

# CONFIG
TCP_IP = '127.0.0.1'
TCP_PORT = 3333
AI_NAME = "Python AI"
AI_COLOR = "orange"

# ------ global vars and data ---------------------------------------------------------------------------------- #

# thread lock
lock = Lock()

# STATUS
Iteration = 0
Endtime = 0
MaxUpdateTime = "0s"
MaxPlayers = 0

# PLAYERS
PlayerID = -1
Name = [None]*1000
Color = [None]*1000
Position = [None]*1000
Velocity = [None]*1000
Acceleration = [None]*1000
Score = [None]*1000
Angle = [None]*1000
TouchingCells = [None]*1000
IsAlive = [None]*1000

# MAP
Grid = []

# ------ start connection -------------------------------------------------------------------------------------------- #


# TCP connection
conn = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
conn.connect((TCP_IP, TCP_PORT))

# init game
cmd = "myAiClient|" + AI_NAME + "|" + AI_COLOR
conn.send(bytes(cmd, 'utf8') + b'\n')
print("TO SERVER:   ", cmd)
resp = conn.makefile().readline()
print("FROM SERVER: ", resp)

# set player id or exit with error
pa = resp.split(":")
if len(pa) == 2:
    PlayerID = int(pa[1].strip())
else:
    print("ERROR!")
    exit(1)

# ------ read & update world-status ---------------------------------------------------------------------------------- #


# parse STATUS block
def parse_stat(text):
    lock.acquire()  # <---- LOCK
    for line in text.split("\n"):
        args = line.split(":")
        if len(args) == 2:
            if args[0] == "Iteration":
                global Iteration
                Iteration = int(args[1].strip())
            elif args[0] == "Endtime":
                global Endtime
                Endtime = int(args[1].strip())
            elif args[0] == "MaxUpdateTime":
                global MaxUpdateTime
                MaxUpdateTime = args[1].strip()
            elif args[0] == "MaxPlayers":
                global MaxPlayers
                MaxPlayers = int(args[1].strip())
    lock.release()  # <---- UNLOCK


# parse PLAYER block
def parse_ply(text):
    lock.acquire()  # <---- LOCK
    current_player_id = -1
    for line in text.split("\n"):
        # split elements
        els = line.split("|")
        if len(els) == 10:
            for el in els:
                # split args
                args = el.split(":")
                if len(args) == 2:
                    # process args
                    if args[0] == "PlayerID":
                        current_player_id = int(args[1].strip())
                    elif args[0] == "Name":
                        global Name
                        Name[current_player_id] = args[1].strip()
                    elif args[0] == "Color":
                        global Color
                        Color[current_player_id] = args[1].strip()
                    elif args[0] == "Position":
                        global Position
                        Position[current_player_id] = args[1].strip()
                    elif args[0] == "Velocity":
                        global Velocity
                        Velocity[current_player_id] = args[1].strip()
                    elif args[0] == "Acceleration":
                        global Acceleration
                        Acceleration[current_player_id] = args[1].strip()
                    elif args[0] == "Score":
                        global Score
                        Score[current_player_id] = args[1].strip()
                    elif args[0] == "Angle":
                        global Angle
                        Angle[current_player_id] = args[1].strip()
                    elif args[0] == "TouchingCells":
                        global TouchingCells
                        TouchingCells[current_player_id] = args[1].strip()
                    elif args[0] == "IsAlive":
                        global IsAlive
                        IsAlive[current_player_id] = args[1].strip()
    lock.release()  # <---- UNLOCK


# parse MAP block
def parse_map(text):
    lock.acquire()  # <---- LOCK
    # reset grid
    global Grid
    Grid = []
    # update grid
    lines = text.split("\n")
    for yH in range(len(lines)):
        line = lines[yH]
        if len(line) > 0:
            for xW in range(len(line)):
                if len(Grid) <= xW:
                    Grid.append([])
                if len(Grid[xW]) <= yH:
                    Grid[xW].append([])
                Grid[xW][yH] = line[xW]
    # fin
    lock.release()  # <---- UNLOCK


# status_updater receives status updates continuously and updates the global variables
def status_updater() -> None:
    s_stat = "START STATUS"
    e_stat = "END STATUS"
    s_ply = "START PLAYER"
    e_ply = "END PLAYER"
    s_map = "START MAP"
    e_map = "END MAP"
    text = ""

    fh = conn.makefile()
    while True:
        line = fh.readline()

        # start new block (reset old block)
        if line.startswith(s_stat) or line.startswith(s_ply) or line.startswith(s_map):
            text = ""
            continue

        # end current block
        if line.startswith(e_stat):
            parse_stat(text)
            continue
        if line.startswith(e_ply):
            parse_ply(text)
            continue
        if line.startswith(e_map):
            parse_map(text)
            continue

        # process block
        text += line


# run status_updater in the background
thread = Thread(target=status_updater)
thread.start()

# wait for server data
while True:
    time.sleep(0.1)
    lock.acquire()  # <---- LOCK
    if Name[0] is not None and len(Grid) > 0:
        lock.release()  # <---- UNLOCK 1/2
        break
    lock.release()  # <---- UNLOCK 2/2

print("go go go ...")


# send cmd function
def set_move_cmd(vx, vy):
    mv = str(vx) + "|" + str(vy)
    conn.send(bytes(mv, 'utf8') + b'\n')
    print("SET MOVE: ", mv)


# ----- implement your AI here --------------------------------------------------------------------------------------- #

while True:
    lock.acquire()  # <---- LOCK (don't remove this!)
    ########################################################

    # Example: print all WORLD STATUS
    print("Iteration", Iteration)
    print("Endtime", Endtime)
    print("MaxUpdateTime", MaxUpdateTime)
    print("MaxPlayers", MaxPlayers)

    # Example: print all PLAYER data
    print("PlayerID", PlayerID)
    print("Name", Name[PlayerID])
    print("Color", Color[PlayerID])
    print("Position", Position[PlayerID])
    print("Velocity", Velocity[PlayerID])
    print("Acceleration", Acceleration[PlayerID])
    print("Score", Score[PlayerID])
    print("Angle", Angle[PlayerID])
    print("TouchingCells", TouchingCells[PlayerID])
    print("IsAlive", IsAlive[PlayerID])

    # Example: set move command
    # don't send the command with every tick!
    if Iteration == 60:
        set_move_cmd(120.3, 300.123)
    if Iteration == 180:
        set_move_cmd(-120.3, 300.123)
    if Iteration == 250:
        set_move_cmd(-120.3, -300.123)

    # Example: describe the fields that i am currently touching
    for cell in TouchingCells[PlayerID].split(";"):
        # get x and y
        xy = cell.split(",")
        if len(xy) == 2 and xy[0] != "" and xy[1] != "":
            x = int(xy[0])
            y = int(xy[1])
            # check: out of range
            if len(Grid) > x and len(Grid[x]) > y:
                # print cell
                typ = Grid[x][y]
                print(" > cell " + str(xy) + " is '" + str(typ) + "'")

    print("--------------------------------")

    ########################################################
    lock.release()  # <---- UNLOCK (don't remove this!)
    time.sleep(0.010)  # > 60 ticks per sec
