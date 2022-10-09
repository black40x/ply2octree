# Ply2Octree

Convert Ply file to octree chunked.


```
go build -o ply2octree
```

```
$> ply2octree <ply-file> <otput-dir>
```

Output binary contains data:
- X,Y,Z float64 [24] byte
- R,G,B [3] byte

meta.json contain bounding box, hierarchy, spacing
