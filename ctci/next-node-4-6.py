#!/usr/bin/python

import sys

class TreeNode:
    def __init__(self, value):
        self.left = None
        self.right = None
        self.parent = None
        self.value = value

    def insert(self, node):
        if node.value > self.value:
            if self.right == None:
                node.parent = self
                self.right = node
            else:
                self.right.insert(node)
        else:
            if self.left == None:
                node.parent = self
                self.left = node
            else:
                self.left.insert(node)

def min_node(compare_node, *nodes):
    lowest_node = TreeNode(sys.maxint)
    for node in nodes:
        if node.value < lowest_node.value and node.value > compare_node.value:
            lowest_node = node
    return lowest_node

def find_next_in_whole_tree(node):
    compare_node = node
    while node.parent != None:
        node = node.parent

    return find_next_in_subtree(node, compare_node)

def find_next_in_subtree(cur_node, compare_node):
    max_node = TreeNode(sys.maxint)

    if cur_node == None:
        return max_node

    lowest_left = max_node

    # Only go left if not our compare node
    if cur_node != compare_node:
        lowest_left = find_next_in_subtree(cur_node.left, compare_node)

    lowest_right = find_next_in_subtree(cur_node.right, compare_node)

    return min_node(compare_node, cur_node, lowest_left, lowest_right)


root = TreeNode(100)
a = TreeNode(50)
b = TreeNode(700)

root.insert(a)
root.insert(b)
root.insert(TreeNode(234123))
root.insert(TreeNode(23))
root.insert(TreeNode(51))

print find_next_in_whole_tree(a).value
