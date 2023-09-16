<template>

  <div class="home">
    <nav>
      <ul>
        <li><a href="/#/login">登录</a></li>
        <li><a href="/#/register">注册</a></li>
      </ul>
    </nav>
    <br>
    <h1 style="border-bottom: 2px dashed #d8d8d8; margin-bottom: 1.4rem;">Welcome to 新星河生态数据web3.0社区</h1>

    <h2>一个可以构建、分享、交易生态数据的社区平台</h2>

    <main>
      <div class="product-list">
        <div v-for="product in products" :key="product.id" class="product-card">
          <img :src="product.image" alt="Product Image">
          <h3>{{ product.name }}</h3>
          <p>{{ product.description }}</p>

        </div>
      </div>
    </main>


  </div>
</template>

<script>
export default {
  data() {
    return {
      products: [
        {
          id: 1,
          name: '基因测序',
          description: '如生物学及生物技术行业功能基因数据、生态环境科学的宏基因组数据',
          image: '1.jpg'
        },
        {
          id: 2,
          name: '环境检测',
          description: '环境污染的治理及监控数据',
          image: '2.jpg'
        },
        {
          id: 3,
          name: '家庭社交',
          description: '为会员进行生育健康咨询、基因测序服务',
          image: '3.jpg'
        },
        {
          id: 4,
          name: '区块链',
          description: '基于区块链技术的个人基因数据保护',
          image: '4.png'
        }
      ],
      
    };
  },
  methods: {

  }
};
</script>

<style scoped>
.home {
  position: relative;
  text-align: center;
}

h1 {
  display: block;
  margin-top: 20px;
}

nav {
  position: relative;
}

nav ul {
  float: right;
  margin-right: 20px;
}

nav ul li {
  float: left;

  margin-left: 10px;
  list-style: none;
}

nav ul li a {
  text-decoration: none;
}

h2 {
  position: absolute;
  right: 50px;
}

main {
  margin: 20px;
  margin-top: 100px;
}

.product-list {
  display: flex;
  flex-wrap: wrap;
  /* justify-content: center;
  align-items: flex-start; */
}

.product-card {
  width: 800px;
  height: 500px;
  margin: 15px;

  padding: 10px;
  border: 1px solid #ccc;
}

.product-card {
  width: 45%;
}

.product-card img {
  width: 100%;
  max-height: 350px;
  object-fit: cover;
  margin-bottom: 10px;
}

.success-message {
  background-color: #4CAF50;
  color: white;
  padding: 10px;
}
</style>
