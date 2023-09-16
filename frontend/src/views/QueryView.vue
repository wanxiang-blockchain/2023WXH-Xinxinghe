<script setup>

</script>

<template>
  <el-container>
    <el-header class="header">
      <h1 class="title">基因数据</h1>
      <input type="search" class="search" placeholder="输入您感兴趣的标签">
    </el-header>
    <el-container>
      <el-aside>
        <el-tag class="tag">标签一</el-tag>
        <el-tag class="tag" type="success">标签二</el-tag>
        <el-tag class="tag" type="info">标签三</el-tag>
        <el-tag class="tag" type="warning">标签四</el-tag>
        <el-tag class="tag" type="danger">标签五</el-tag>
      </el-aside>
      <el-main>Main</el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.header {
  border-bottom: 2px dashed #d8d8d8;
  margin-bottom: 1.4rem;
}

.title {
  text-align: left;
  display: inline-block;
  float: left;
}

.search {
  padding: 0.5rem 20rem 0.5rem 2.5rem;
  margin-top: 0.2rem;
}

.tag {
  float: left;
  display: block;
}
</style>