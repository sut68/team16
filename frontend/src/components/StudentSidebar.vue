<template>
  <div :class="{ 'sidebar-closed': !isSidebarOpen, 'sidebar-mobile-open': isMobileMenuOpen }">
    <!-- Mobile Header Bar -->
    <div class="mobile-header">
      <button @click="toggleMobileMenu" class="mobile-menu-btn" aria-label="เปิดเมนู">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5" />
        </svg>
      </button>
      <div class="mobile-header-title">
        <img class="mobile-logo" :src="logo" alt="logo">
        <span>ทุนการศึกษา</span>
      </div>
      <div class="mobile-header-spacer"></div>
    </div>

    <!-- Mobile Overlay -->
    <div 
      v-if="isMobileMenuOpen" 
      class="mobile-overlay"
      @click="closeMobileMenu"
    ></div>

    <aside ref="asideRef" :class="{ 'mobile-open': isMobileMenuOpen }">
      <!-- Mobile Close Button -->
      <button @click="closeMobileMenu" class="mobile-close-btn" aria-label="ปิดเมนู">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <div class="logo-and-menu-header">
        <img class="logo" :src="logo" alt="logo">
        <div class="header-text-wrapper" v-show="isSidebarOpen || isMobileMenuOpen">
          <span> ทุนการศึกษา </span>
          <span> สำนักวิชาวิศวกรรมศาสตร์ </span>
        </div>
      </div>

      <router-link to="/dashboard" class="menu-link" exact-active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M3.75 6A2.25 2.25 0 0 1 6 3.75h2.25A2.25 2.25 0 0 1 10.5 6v2.25a2.25 2.25 0 0 1-2.25 2.25H6a2.25 2.25 0 0 1-2.25-2.25V6ZM3.75 15.75A2.25 2.25 0 0 1 6 13.5h2.25a2.25 2.25 0 0 1 2.25 2.25V18a2.25 2.25 0 0 1-2.25 2.25H6A2.25 2.25 0 0 1 3.75 18v-2.25ZM13.5 6a2.25 2.25 0 0 1 2.25-2.25H18A2.25 2.25 0 0 1 20.25 6v2.25A2.25 2.25 0 0 1 18 10.5h-2.25a2.25 2.25 0 0 1-2.25-2.25V6ZM13.5 15.75a2.25 2.25 0 0 1 2.25-2.25H18a2.25 2.25 0 0 1 2.25 2.25V18A2.25 2.25 0 0 1 18 20.25h-2.25A2.25 2.25 0 0 1 13.5 18v-2.25Z" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">แดชบอร์ด</span>
      </router-link>

      <router-link to="/dashboard/apply-scholarship" class="menu-link" active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">ยื่นสมัครทุน</span>
      </router-link>

      <router-link to="/dashboard/track-status" class="menu-link" active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607zM10.5 7.5v6l3.75 2.25" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">ติดตามสถานะทุน</span>
      </router-link>

      <router-link to="/dashboard/schedule" class="menu-link" active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607zM10.5 7.5v6l3.75 2.25" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">จองรอบสัมภาษณ์</span>
      </router-link>

      <router-link to="/dashboard/assistance" class="menu-link" active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M9.879 7.519c1.171-1.025 3.071-1.025 4.242 0 1.172 1.025 1.172 2.687 0 3.712-.203.179-.43.326-.67.442-.745.361-1.45.999-1.45 1.827v.75M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 5.25h.008v.008H12v-.008Z" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">ติดต่อช่วยเหลือ</span>
      </router-link>

      <router-link to="/dashboard/profilestudent" class="menu-link" active-class="active" @click="closeMobileMenu">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5"
          class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round"
            d="M17.982 18.725A7.488 7.488 0 0 0 12 15.75a7.488 7.488 0 0 0-5.982 2.975m11.963 0a9 9 0 1 0-11.963 0m11.963 0A8.966 8.966 0 0 1 12 21a8.966 8.966 0 0 1-5.982-2.275M15 9.75a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
        </svg>

        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">โปรไฟล์</span>
      </router-link>

      <div class="flex-grow"></div> <!-- Spacer -->

      <a @click="handleLogout" class="menu-link cursor-pointer mt-auto mb-4">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0 0 13.5 3h-6a2.25 2.25 0 0 0-2.25 2.25v13.5A2.25 2.25 0 0 0 7.5 21h6a2.25 2.25 0 0 0 2.25-2.25V15m3 0 3-3m0 0-3-3m3 3H9" />
        </svg>
        <span v-show="isSidebarOpen || isMobileMenuOpen" class="menu-text">ออกจากระบบ</span>
      </a>

    </aside>
    <main class="main-content">
      <button @click="toggleSidebar" class="sidebar-toggle-btn" :style="{ top: toggleButtonTop }">
        <svg :class="{ 'rotate-180': !isSidebarOpen }" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5L8.25 12l7.5-7.5" />
        </svg>
      </button>
      <div class='main-scroll'>
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
    const isMobileMenuOpen = ref(false);
    const toggleButtonTop = ref('180px');
    const asideRef = ref<HTMLElement | null>(null);
    const route = useRoute();
    const router = useRouter();
    
    const TRANSITION_DURATION = 300; // ms
    
    // สำหรับ debounce
    let updateTimeout: ReturnType<typeof setTimeout> | null = null;
    let resizeObserver: ResizeObserver | null = null;

    const handleLogout = () => {
      try { sessionStorage.clear(); } catch {}
      try { localStorage.clear(); } catch {}
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

    // Mobile menu functions
    const toggleMobileMenu = () => {
      isMobileMenuOpen.value = !isMobileMenuOpen.value;
      // Prevent body scroll when mobile menu is open
      if (isMobileMenuOpen.value) {
        document.body.style.overflow = 'hidden';
      } else {
        document.body.style.overflow = '';
      }
    };

    const closeMobileMenu = () => {
      isMobileMenuOpen.value = false;
      document.body.style.overflow = '';
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
      isMobileMenuOpen,
      toggleSidebar,
      toggleMobileMenu,
      closeMobileMenu,
      toggleButtonTop,
      asideRef,
      handleLogout,
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
  margin-top: 2px;
  margin-bottom: 2px;
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
  width: 50px;
  margin-top: 25px;
  margin-bottom: 10px;
  margin-left: -5px;
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
  margin-top: -30px;
  padding: 20px 0;
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
  bottom: 5px;
  left: -15px;
  width: 95%;
  height: 2.5px;
  
  background-color: #ffffff;
  margin-left: 10px;
}

div.sidebar-closed .logo-and-menu-header::after {
  display: none;
}

/* ========================================
   MOBILE RESPONSIVE STYLES
   ======================================== */

/* Mobile Header Bar - visible only on mobile */
.mobile-header {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: linear-gradient(135deg, #8B001D 0%, #6d0016 100%);
  z-index: 1000;
  padding: 0 16px;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.mobile-menu-btn {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  padding: 10px;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.mobile-menu-btn:hover,
.mobile-menu-btn:active {
  background: rgba(255, 255, 255, 0.2);
}

.mobile-header-title {
  display: flex;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 16px;
  font-weight: 600;
}

.mobile-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.mobile-header-spacer {
  width: 44px; /* Same as button for centering */
}

/* Mobile Overlay */
.mobile-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 998;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Mobile Close Button inside sidebar */
.mobile-close-btn {
  display: none;
  position: absolute;
  top: 16px;
  right: 16px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  padding: 8px;
  border-radius: 50%;
  cursor: pointer;
  z-index: 10;
  transition: background 0.2s;
}

.mobile-close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* ========================================
   MOBILE BREAKPOINT (max-width: 768px)
   ======================================== */
@media (max-width: 768px) {
  /* Show mobile header */
  .mobile-header {
    display: flex;
  }

  /* Show overlay when menu is open */
  .mobile-overlay {
    display: block;
  }

  /* Show close button in sidebar */
  .mobile-close-btn {
    display: block;
  }

  /* Hide desktop toggle button on mobile */
  .sidebar-toggle-btn {
    display: none !important;
  }

  /* Sidebar - hidden by default, slide in from left */
  aside {
    position: fixed;
    left: -280px; /* Start off-screen */
    top: 0;
    width: 280px;
    height: 100vh;
    z-index: 999;
    border-radius: 0;
    border-top-right-radius: 30px;
    border-bottom-right-radius: 30px;
    padding-top: 20px;
    padding-left: 15px;
    transition: left 0.3s ease-in-out;
  }

  /* Sidebar open state */
  aside.mobile-open {
    left: 0;
  }

  /* Main content - full width on mobile */
  .main-content {
    margin-left: 0;
    margin-top: 60px; /* Space for mobile header */
    height: calc(100vh - 60px);
  }

  /* Main scroll area adjustments */
  .main-scroll {
    border-top-left-radius: 0;
    margin: 0;
    padding: 12px;
  }

  /* Logo adjustments */
  .logo-and-menu-header {
    margin-top: 10px;
    padding-right: 50px; /* Space for close button */
  }

  /* Menu link touch-friendly sizing */
  .menu-link {
    padding: 14px 12px 14px 20px;
    font-size: 15px;
    min-height: 48px; /* Touch-friendly minimum */
  }

  /* Override sidebar-closed styles on mobile */
  div.sidebar-closed aside {
    width: 280px;
    left: -280px;
  }

  div.sidebar-closed .main-content {
    margin-left: 0;
  }

  div.sidebar-closed .menu-link {
    padding-left: 20px;
    justify-content: flex-start;
    border-radius: 0;
    margin-left: 0;
    margin-right: 0;
  }

  div.sidebar-closed .logo {
    margin: 0;
    margin-top: 25px;
    margin-bottom: 10px;
    margin-left: -5px;
    width: 50px;
  }

  div.sidebar-closed .logo-and-menu-header {
    justify-content: flex-start;
    padding-left: inherit;
  }

  div.sidebar-closed .menu-link svg {
    margin-right: 0 !important;
  }
}

/* ========================================
   SMALL MOBILE (max-width: 480px)
   ======================================== */
@media (max-width: 480px) {
  aside {
    width: 100%;
    left: -100%;
    border-radius: 0;
  }

  aside.mobile-open {
    left: 0;
  }

  div.sidebar-closed aside {
    width: 100%;
    left: -100%;
  }

  .mobile-header-title span {
    font-size: 14px;
  }

  .main-scroll {
    padding: 8px;
  }
}

</style>

