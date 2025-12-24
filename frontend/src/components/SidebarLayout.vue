<template>
  <div :class="{ 'sidebar-closed': !isSidebarOpen }" class="h-screen" data-theme="light" data-testid="sidebar-layout-container">
    <aside ref="asideRef" data-testid="sidebar-aside">
      <div class="logo-and-menu-header">
        <img class="logo" :src="logo" alt="logo" data-testid="sidebar-logo">
        <div class="header-text-wrapper" v-show="isSidebarOpen">
          <span> ทุนการศึกษา </span>
          <span> สำนักวิชาวิศวกรรมศาสตร์ </span>
        </div>
      </div>

      <router-link to="/admin" class="menu-link" exact-active-class="active" data-testid="sidebar-menu-dashboard">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z" />
        </svg>
        <span v-show="isSidebarOpen" class="menu-text">แดชบอร์ด</span>
      </router-link>

      <router-link to="/admin/news" class="menu-link" :class="{ active: isMenuActive('/admin/news') }" data-testid="sidebar-menu-news">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M10.34 15.84c-.688-.06-1.386-.09-2.09-.09H7.5a4.5 4.5 0 1 1 0-9h.75c.704 0 1.402-.03 2.09-.09m0 9.18c.253.962.584 1.892.985 2.783.247.55.06 1.21-.463 1.511l-.657.38c-.551.318-1.26.117-1.527-.461a20.845 20.845 0 0 1-1.44-4.282m3.102.069a18.03 18.03 0 0 1-.59-4.59c0-1.586.205-3.124.59-4.59m0 9.18a23.848 23.848 0 0 1 8.835 2.535M10.34 6.66a23.847 23.847 0 0 0 8.835-2.535m0 0A23.74 23.74 0 0 0 18.795 3m.38 1.125a23.91 23.91 0 0 1 1.014 5.395m-1.014 8.855c-.118.38-.245.754-.38 1.125m.38-1.125a23.91 23.91 0 0 0 1.014-5.395m0-3.46c.495.413.811 1.035.811 1.73 0 .695-.316 1.317-.811 1.73m0-3.46a24.347 24.347 0 0 1 0 3.46" />
        </svg>
        <span v-show="isSidebarOpen" class="menu-text">ข่าวสาร</span>
      </router-link>

      <router-link to="/admin/sponsors" class="menu-link" :class="{ active: isMenuActive('/admin/sponsors') }" data-testid="sidebar-menu-sponsors">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M2.25 21h19.5m-18-18v18m10.5-18v18m6-13.5V21M6.75 6.75h.75m-.75 3h.75m-.75 3h.75m3-6h.75m-.75 3h.75m-.75 3h.75M6.75 21v-3.375c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21M3 3h12m-.75 4.5H21m-3.75 3.75h.008v.008h-.008v-.008Zm0 3h.008v.008h-.008v-.008Zm0 3h.008v.008h-.008v-.008Z" />
        </svg>
        <span v-show="isSidebarOpen" class="menu-text">บริษัทที่ให้ทุน</span>
      </router-link>

      <router-link to="/admin/projects" class="menu-link" :class="{ active: isMenuActive('/admin/projects') }" data-testid="sidebar-menu-projects">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 0 0 2.25-2.25V6.75A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25v10.5A2.25 2.25 0 0 0 4.5 19.5Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">โครงการทุน</span>
      </router-link>

      <router-link to="/admin/screening" class="menu-link" :class="{ active: isMenuActive('/admin/screening') }" data-testid="sidebar-menu-screening">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M15 9h3.75M15 12h3.75M15 15h3.75M4.5 19.5h15a2.25 2.25 0 0 0 2.25-2.25V6.75A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25v10.5A2.25 2.25 0 0 0 4.5 19.5Zm6-10.125a1.875 1.875 0 1 1-3.75 0 1.875 1.875 0 0 1 3.75 0Zm1.294 6.336a6.721 6.721 0 0 1-3.17.789 6.721 6.721 0 0 1-3.168-.789 3.376 3.376 0 0 1 6.338 0Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">คัดกรอง</span>
      </router-link>

      <router-link to="/admin/approval" class="menu-link" :class="{ active: isMenuActive('/admin/approval') }" data-testid="sidebar-menu-approval">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M10.125 2.25h-4.5c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125v-9M10.125 2.25h.375a9 9 0 0 1 9 9v.375M10.125 2.25A3.375 3.375 0 0 1 13.5 5.625v1.5c0 .621.504 1.125 1.125 1.125h1.5a3.375 3.375 0 0 1 3.375 3.375M9 15l2.25 2.25L15 12" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">อนุมัติเอกสาร</span>
      </router-link>

      <router-link to="/admin/interview" class="menu-link" :class="{ active: isMenuActive('/admin/interview') }" data-testid="sidebar-menu-interview">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 0 1 2.25-2.25h13.5A2.25 2.25 0 0 1 21 7.5v11.25m-18 0A2.25 2.25 0 0 0 5.25 21h13.5A2.25 2.25 0 0 0 21 18.75m-18 0v-7.5A2.25 2.25 0 0 1 5.25 9h13.5A2.25 2.25 0 0 1 21 11.25v7.5m-9-6h.008v.008H12v-.008ZM12 15h.008v.008H12V15Zm0 2.25h.008v.008H12v-.008ZM9.75 15h.008v.008H9.75V15Zm0 2.25h.008v.008H9.75v-.008ZM7.5 15h.008v.008H7.5V15Zm0 2.25h.008v.008H7.5v-.008Zm6.75-4.5h.008v.008h-.008v-.008Zm0 2.25h.008v.008h-.008V15Zm0 2.25h.008v.008h-.008v-.008Zm2.25-4.5h.008v.008H16.5v-.008Zm0 2.25h.008v.008H16.5V15Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">จัดการการสัมภาษณ์</span>
      </router-link>

      <router-link to="/admin/evaluations" class="menu-link" :class="{ active: isMenuActive('/admin/evaluations') }" data-testid="sidebar-menu-consider">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">พิจารณา</span>
      </router-link>

      <router-link to="/admin/assistance" class="menu-link" :class="{ active: isMenuActive('/admin/assistance') }" data-testid="sidebar-menu-assistance">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">ติดต่อช่วยเหลือ</span>
      </router-link>

      <router-link to="/admin/users" class="menu-link" :class="{ active: isMenuActive('/admin/users') }" data-testid="sidebar-menu-users">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path strokeLinecap="round" strokeLinejoin="round"
            d="M18 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 1 1-6.75 0 3.375 3.375 0 0 1 6.75 0ZM3 19.235v-.11a6.375 6.375 0 0 1 12.75 0v.109A12.318 12.318 0 0 1 9.374 21c-2.331 0-4.512-.645-6.374-1.766Z" />
        </svg>


        <span v-show="isSidebarOpen" class="menu-text">จัดการผู้ใช้</span>
      </router-link>

      <router-link to="/admin/profileadmin" class="menu-link" :class="{ active: isMenuActive('/admin/profileadmin') }" data-testid="sidebar-menu-profile">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
        </svg>

        <span v-show="isSidebarOpen" class="menu-text">โปรไฟล์</span>
      </router-link>

      <div class="flex-grow"></div> <!-- Spacer -->

      <a @click="handleLogout" class="menu-link cursor-pointer mt-auto mb-4" data-testid="sidebar-menu-logout">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
        </svg>
        <span v-show="isSidebarOpen" class="menu-text">ออกจากระบบ</span>
      </a>

    </aside>
    <main class="main-content h-screen flex flex-col">
      <button @click="toggleSidebar" class="sidebar-toggle-btn" :style="{ top: toggleButtonTop }" data-testid="sidebar-toggle-button">
        <svg :class="{ 'rotate-180': !isSidebarOpen }" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
        </svg>
      </button>
      <div class='main-scroll flex-1 overflow-auto' data-testid="sidebar-main-content">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, nextTick, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import logo from '@/assets/logo/ENGi Logo-White.png';

export default defineComponent({
  name: 'SidebarMenu',
  setup() {
    const isSidebarOpen = ref(true);
    const toggleButtonTop = ref('180px');
    const asideRef = ref<HTMLElement | null>(null);
    const route = useRoute();
    const router = useRouter();

    const isMenuActive = (basePath: string): boolean => {
      return route.path === basePath || route.path.startsWith(basePath + '/');
    }

    const TRANSITION_DURATION = 300; // ms

    // สำหรับ debounce
    let updateTimeout: ReturnType<typeof setTimeout> | null = null;
    let resizeObserver: ResizeObserver | null = null;

    const handleLogout = () => {
      try { sessionStorage.clear(); } catch { }
      try { localStorage.clear(); } catch { }
      router.push('/');
    };

    const updateButtonPosition = () => {
      if (updateTimeout) {
        clearTimeout(updateTimeout);
      }

      updateTimeout = setTimeout(async () => {
        await nextTick();

        if (!asideRef.value) return;

        const activeLink = asideRef.value.querySelector('.menu-link.active') as HTMLElement;

        if (activeLink) {
          const asideRect = asideRef.value.getBoundingClientRect();
          const linkRect = activeLink.getBoundingClientRect();

          const buttonHeight = 40;
          const linkCenterY = linkRect.top - asideRect.top + (linkRect.height / 2);
          const newTop = Math.max(0, linkCenterY - (buttonHeight / 2));

          toggleButtonTop.value = `${newTop}px`;
        }
      }, 50); // debounce 50ms
    };

    const toggleSidebar = () => {
      isSidebarOpen.value = !isSidebarOpen.value;

      // รอให้ transition เสร็จก่อนคำนวณใหม่
      setTimeout(() => {
        updateButtonPosition();
      }, TRANSITION_DURATION);
    };

    const setupResizeObserver = () => {
      if (!asideRef.value) return;

      resizeObserver = new ResizeObserver(() => {
        updateButtonPosition();
      });

      resizeObserver.observe(asideRef.value);
    };

    onMounted(() => {

      nextTick(() => {
        updateButtonPosition();
        setupResizeObserver();
      });
    });

    onBeforeUnmount(() => {
      // cleanup
      if (updateTimeout) {
        clearTimeout(updateTimeout);
      }
      if (resizeObserver) {
        resizeObserver.disconnect();
      }
    });

    // watch route changes
    watch(
      () => route.path,
      () => {
        updateButtonPosition();
      }
    );

    return {
      logo,
      isSidebarOpen,
      toggleSidebar,
      toggleButtonTop,
      asideRef,
      handleLogout,
      isMenuActive,
    };
  },
});
</script>

<style scoped>
aside,
.main-content {
  transition: all 0.3s ease-in-out;
}

.menu-link {
  font-size: 15px;
  color: #fff;
  display: flex;
  align-items: center;
  padding: 12px;
  padding-left: 30px;
  text-decoration: none;
  -webkit-tap-highlight-color: transparent;
  margin-top: 5px;
  margin-bottom: 5px;
  position: relative;
}

.menu-text {
  flex: 1;
  margin-left: 10px;
}

.sidebar-toggle-btn {
  display: flex;
  background: #253C90;
  color: white;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  z-index: 99;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  position: absolute;
  left: -30px;
  transition: top 0.3s ease-in-out, transform 0.3s ease-in-out, background-color 0.3s ease;
}

.sidebar-toggle-btn:hover {
  background-color: #3352c3;
  transform: scale(1.1);
}

.sidebar-toggle-btn svg {
  width: 20px;
  height: 20px;
  transition: transform 0.3s ease-in-out;
}

.sidebar-toggle-btn svg.rotate-180 {
  transform: rotate(180deg);
}

.menu-link:hover,
.menu-link.active {
  color: #000000;
  background: #fff;
  outline: none;
  background-color: #fff;
  border-top-left-radius: 20px;
  border-bottom-left-radius: 20px;
}

.menu-link:hover::after,
.menu-link.active::after {
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

.menu-link:hover::before,
.menu-link.active::before {
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

div.sidebar-closed aside {
  width: 80px;
  padding-left: 0;
  border-top-right-radius: 50px;
}

div.sidebar-closed .main-content {
  margin-left: 90px;
}

div.sidebar-closed .logo-and-menu-header {
  justify-content: center;
  padding-left: 0;
}

div.sidebar-closed .logo {
  margin: 0;
  margin-top: 30px;
  width: 50px;
  margin-bottom: -20px;
}

div.sidebar-closed .menu-link {
  padding-left: 0;
  justify-content: center;
  border-radius: 15px;
  margin-left: 10px;
  margin-right: 10px;
  margin-top: 10px;
}

div.sidebar-closed .menu-link svg {
  margin-right: -8px !important;
}

div.sidebar-closed .menu-link:hover,
div.sidebar-closed .menu-link.active {
  border-radius: 15px;
}

div.sidebar-closed .menu-link:hover::after,
div.sidebar-closed .menu-link.active::after,
div.sidebar-closed .menu-link:hover::before,
div.sidebar-closed .menu-link.active::before {
  display: none;
}

aside {
  color: #fff;
  width: 250px;
  padding-left: 20px;
  height: 100vh;
  background-image: linear-gradient(30deg, #8B001D, #8B001D);
  border-top-right-radius: 130px;
  position: absolute;
  left: 0;
  overflow: hidden;
}

aside .logo {
  width: 60px;
  margin-top: 35px;
  margin-bottom: 20px;
  margin-left: -10px;
  flex-shrink: 0;
  transition: all 0.3s ease-in-out;
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
  transition: all 0.3s ease-in-out;

}

.header-text-wrapper {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
  transform: translateY(10px);
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
  height: 100vh;
  position: relative;
  display: flex;
  flex-direction: column;
  /* overflow: visible; */
  margin-left: 260px;
}

.main-scroll {
  flex: 1;
  width: 100%;
  padding: 20px;
  /* height: 100%; */
  border-top-left-radius: 50px;
  background-color: #f0f2f5;
  overflow-x: hidden;
  overflow-y: hidden;
  margin-left: 30px;
  /* max-width: 1500px; */
  margin: 0 0;
}

.logo-and-menu-header {
  position: relative;
  padding-bottom: 20px;
}

.logo-and-menu-header::after {
  content: "";
  position: absolute;
  bottom: 17px;
  left: -15px;
  width: 95%;
  height: 2.5px;
  background-color: #ffffff;
  margin-left: 10px;
}

div.sidebar-closed .logo-and-menu-header::after {
  display: none;
}
</style>
