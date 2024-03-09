#!/bin/sh
cat .devcontainer/bashrc.override.sh >> ~/.bashrc

# go get -u golang.org/x/tools/gopls
# go get -u github.com/ramya-rao-a/go-outline
# go get -u github.com/acroca/go-symbols
# # 此外，为了更完整的开发体验，您可能还需要安装其他相关工具：
# go get -u github.com/uudashr/gopkgs/v2/cmd/gopkgs  # 包信息查询
# go get -u github.com/rogpeppe/godef          # 跳转到定义
# go get -u honnef.co/go/tools/cmd/staticcheck # 高级静态检查工具

# # 如果需要生成或更新依赖关系图
# go get -u github.com/kisielk/godepgraph     # 生成依赖关系图
