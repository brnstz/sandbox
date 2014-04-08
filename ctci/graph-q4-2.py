#!/usr/bin/python

class node: 
    def __init__(self, value):
        self.vertices = []
        self.value = value

    def add_vertex(self, node):
        self.vertices.append(node)

def check_route_from(node_orig, node_dest, node_loopcheck = {}):
    if node_orig == node_dest:
        return True

    if node_loopcheck.get(node_orig):
        return False

    node_loopcheck[node_orig] = True
        
    for vertex in node_orig.vertices:
        if check_route_from(vertex, node_dest, node_loopcheck):
            return True

    return False

def check_both_routes(node1, node2):
    if check_route_from(node1, node2):
        return True
    elif check_route_from(node2, node1):
        return True
    else:
        return False


a = node(1234)
b = node(5678)
c = node(9999)
d = node(888888)
e = node(234234)

a.add_vertex(b)
b.add_vertex(c)
c.add_vertex(e)

print check_both_routes(a, e)

