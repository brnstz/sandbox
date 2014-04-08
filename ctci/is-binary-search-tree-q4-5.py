#!/usr/bin/python

class TreeNode:
    def __init__(self, value):
        self.value = value
        self.left = None
        self.right = None

    def insert(self, node):
        if node.value > self.value:
            if self.right == None:
                self.right = node
            else:
                self.right.insert(node)
        else:
            if self.left == None:
                self.left = node
            else:
                self.left.insert(node)

    def is_binary_search_tree(self):
        left_ok = False
        right_ok = False

        if self.left == None:
            left_ok = True
        else:
            left_ok = self.left.value <= self.value and self.left.is_binary_search_tree()

        if self.right == None:
            right_ok = True
        else:
            right_ok = self.right.value > self.value and self.right.is_binary_search_tree()

        return right_ok and left_ok
            
root = TreeNode(10)
b = TreeNode(100)
c = TreeNode(50)
d = TreeNode(25)
e = TreeNode(2342)

root.insert(b)
root.insert(c)
root.insert(d)
root.insert(e)

print root.is_binary_search_tree()

# This should break binary search
e.right = TreeNode(1)

print root.is_binary_search_tree()

