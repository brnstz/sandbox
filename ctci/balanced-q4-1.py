#!/usr/bin/python

class node:
    left = None
    right = None
    value = None

    def __init__(self, value):
        self.value = value

    def insert(self, node):
        if node.value > self.value and self.right != None:
            self.right.insert(node)
        elif node.value <= self.value and self.left != None:
            self.left.insert(node)
        elif node.value > self.value:
            self.right = node
        else:
            self.left = node



def check_balanced(node):
    if node == None:
        return 0, True
    else:
        left_count, left_balanced = check_balanced(node.left) 
        right_count, right_balanced = check_balanced(node.right)
        new_count = left_count + right_count + 1

        if abs(left_count - right_count) > 1:
            return new_count, False
        else:
            return new_count, left_balanced and right_balanced


root = node(100)
root.insert(node(500))
root.insert(node(50))
root.insert(node(25))
root.insert(node(5000))

count, balanced = check_balanced(root)
print balanced

