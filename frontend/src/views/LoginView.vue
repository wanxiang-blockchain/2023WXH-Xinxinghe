<template>
  <div class="bg">
    <img src="../assets/bg.jpg" style="float: left;display: inline-block;width: 50%;margin-top: 5%">
    <div class="login-container">
      <h2>登录页面</h2>
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="address">钱包地址:</label>
          <input v-model="address" type="text" id="address" placeholder="钱包地址" required>
        </div>
        <div class="form-group">
          <label for="password">密码:</label>
          <input v-model="password" type="password" id="password" placeholder="密码" required>
        </div>
        <button type="submit">登录</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

export default {
  data() {
    return {
      address: '',
      password: ''
    };
  },
  methods: {
    login() {
      // 处理登录逻辑
      if (this.email && this.password) {
        let formData = new FormData();
        formData.append("address", this.address);
        formData.append("password", this.password);
        axios.post("/user/login", formData, {
          headers: {'Content-Type': 'application/x-www-form-urlencoded'},
        }).then(function (resp) {
          switch (resp.status) {
            case 200:
              this.error = "登录成功";
              break;
            case 600:
              let obj = JSON.parse(resp.data);
              this.error = obj.desc;
          }
        }).catch(function (err) {
          console.log(err)
        })
      }
    },
    toRegister() {
      this.$router.push({path: '/Register'})
    }
  }
};
</script>

<style scoped>
.bg {
  background-image: url('../assets/bg.jpg');
  background-repeat: no-repeat;
  background-size: cover;
}

.login-container {
  max-width: 400px;
  padding: 20px;
  background-color: #f5f5f5;
  margin: 10%;
  float: right;
}

.form-group {
  margin-bottom: 16px;
}

label {
  margin-bottom: 8px;
  font-weight: bold;
  float: left;
}

input[type="text"],
input[type="password"] {
  width: 100%;
  padding: 8px;
  font-size: 16px;
  border-radius: 4px;
  border: 1px solid #ccc;
}

button[type="submit"] {
  display: block;
  width: 100%;
  padding: 12px;
  font-size: 16px;
  font-weight: bold;
  color: #fff;
  background-color: #007bff;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>
  