#!/usr/bin/python

# Given a sorted array (increasing) of unique integers, write an algorithm
# to create a binary tree of minimal height.

import collections

class node:
    def __init__(self, value):
        self.value = value
        self.left = None
        self.right = None

    def insert(self, node):
        if node.value > self.value and self.right != None:
            self.right.insert(node)
        elif node.value <= self.value and self.left != None:
            self.right.insert(node)
        elif node.value > self.value:
            self.right = node
        else:
            self.left = node

    def children(self):
        return [self.left, self.right]

        
def print_tree(node, visited={}):
    if node == None:
        return

    if visited.get(node) != None:
        return

    print(node.value)
    visited[node.value] = True

    for child in node.children():
        print_tree(child, visited)

def print_tree_bfs(node):
    visited = {}
    q = collections.deque()
    print(node.value)
    visited[node] = True
    q.append(node)
    
    while len(q) > 0:
        cur_node = q.popleft()
        for child_node in cur_node.children():
            if child_node != None and not visited.get(child_node):
                print(child_node.value), " ",
                q.append(child_node)
        print

        
        


incoming = [1, 2, 3, 4, 5, 6, 7, 8, 9]
len_incoming = len(incoming)
midpoint = len_incoming / 2

root = node(incoming[midpoint])
i = midpoint + 1
j = midpoint - 1

while i < len_incoming and j >= 0:
    root.insert(node(incoming[i]))
    root.insert(node(incoming[j]))

    i += 1
    j -= 1

if j >= 0:
    root.insert(node(incoming[j]))

print_tree(root)
print
print_tree_bfs(root)

