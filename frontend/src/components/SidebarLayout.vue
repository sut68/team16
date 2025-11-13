<template>
  <div>
    <aside>
      <div class="logo-and-menu-header">
        <img class="logo" :src="logo" alt="logo">
        <div class="header-text-wrapper">
          <span> ทุนการศึกษา </span>
          <span> สำนักวิชาวิศวกรรมศาสตร์ </span>
        </div>
      </div>

      <a href="javascript:void(0)">
        <i class="fa fa-user-o" aria-hidden="true"></i>
        แดชบอร์ด
      </a>
      <a href="javascript:void(0)">
        <i class="fa fa-laptop" aria-hidden="true"></i>
        บริษัทที่ให้ทุน
      </a>
      <a href="javascript:void(0)">
        <i class="fa fa-clone" aria-hidden="true"></i>
        โครงการทุน
      </a>
      <a href="javascript:void(0)">
        <i class="fa fa-star-o" aria-hidden="true"></i>
        สมัครทุน
      </a>
      <a href="javascript:void(0)">
        <i class="fa fa-trash-o" aria-hidden="true"></i>
        คัดกรอง
      </a>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import logo from '@/assets/logo/ENGi Logo-White.png';

export default defineComponent({
  name: 'SidebarMenu',
  setup() {
    return {
      logo,
    };
  },
});
</script>

<style scoped>
aside {
  color: #fff;
  width: 250px;
  padding-left: 20px;
  height: 100vh;
  background-image: linear-gradient(30deg, #8B001D, #8B001D);
  border-top-right-radius: 130px;
  position: absolute;
  left: 0;
}

aside .logo {
  width: 70px;
  margin-top: 60px;
  margin-bottom: 20px;
  margin-left: -10px;
  flex-shrink: 0;
}

aside a {
  font-size: 15px;
  color: #fff;
  display: block;
  padding: 12px;
  padding-left: 30px;
  text-decoration: none;
  -webkit-tap-highlight-color: transparent;
  margin-top: -5px;
  margin-bottom: 10px;
}

aside a:hover {
  color: #000000;
  background: #fff;
  outline: none;
  position: relative;
  background-color: #fff;
  border-top-left-radius: 20px;
  border-bottom-left-radius: 20px;
}

aside a i {
  margin-right: 5px;
}

aside a:hover::after {
  content: "";
  position: absolute;
  background-color: transparent;
  bottom: 100%;
  right: 0;
  height: 35px;
  width: 35px;
  border-bottom-right-radius: 20px;
  box-shadow: 0 20px 0 0 #fff;
}

aside a:hover::before {
  content: "";
  position: absolute;
  background-color: transparent;
  top: 46px;
  right: 0;
  height: 35px;
  width: 50px;
  border-top-right-radius: 20px;
  box-shadow: 0 -20px 0 0 #fff;
}

aside p {
  margin-top: 20px;
  padding: 40px 0;
}

.logo-and-menu-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: -40px;
  padding: 40px 0;
}

.header-text-wrapper {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
  transform: translateY(20px);
}

.header-text-wrapper span {
  margin-top: 5px;
  padding: 0;
}

.header-text-wrapper span:first-child {
  font-size: 14px;
  font-weight: 500;
}

.header-text-wrapper span:last-child {
  font-size: 14px;
  font-weight: 500;
}

.main-content {
  margin-left: 265px;
  padding: 20px;
  background-color: #f0f2f5;
  height: 100vh;
  overflow-y: auto;
  border-radius: 50px;
  width: 200vh;
}
</style>