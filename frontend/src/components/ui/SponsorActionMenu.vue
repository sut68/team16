<script setup lang="ts">
  import { ref, nextTick, onMounted, onBeforeUnmount } from "vue";
  import { CheckCircle, Ellipsis, PenLine, Slash, Trash, Users } from "lucide-vue-next";
  import Swal from "sweetalert2";

  const props = defineProps<{
    id: number
    status?: "active" | "inactive" | string;
  }>();

  const emit = defineEmits<{
    edit: [number],
    'edit-contacts': [number],
    'toggle-status': [number],
    delete: [number]
  }>();

  const open = ref(false);
  const btnRef = ref<HTMLElement | null>(null);
  const menuRef = ref<HTMLElement | null>(null);
  const menuStyle = ref<Record<string, string>>({});

  // Global event name for closing other menus
  const CLOSE_ALL_MENUS_EVENT = 'sponsor-menu-close-all';

  function toggleMenu(e: MouseEvent) {
    e.stopPropagation();
    
    if (!open.value) {
      // Dispatch event to close all other menus before opening this one
      window.dispatchEvent(new CustomEvent(CLOSE_ALL_MENUS_EVENT, { detail: props.id }));
    }
    
    open.value = !open.value;
    btnRef.value = e.currentTarget as HTMLElement;
    nextTick(positionMenu);
  }

  function positionMenu() {
    if (!btnRef.value || !menuRef.value) return;

    const btn = btnRef.value.getBoundingClientRect();
    const menu = menuRef.value.getBoundingClientRect();

    // ตำแหน่งคงที่: ใต้ปุ่ม, ชิดขวาของปุ่ม
    let left = btn.right - menu.width;
    let top = btn.bottom + 4;

    // ตรวจสอบขอบจอ (เฉพาะกรณีชนขอบ)
    const pad = 8;
    
    // ถ้าชนขอบซ้าย ให้ชิดซ้ายของปุ่มแทน
    if (left < pad) {
      left = btn.left;
    }
    
    // ถ้าชนขอบล่าง ให้แสดงด้านบนแทน
    if (top + menu.height > window.innerHeight - pad) {
      top = btn.top - menu.height - 4;
    }

    menuStyle.value = {
      position: "fixed",
      left: `${left}px`,
      top: `${top}px`,
      zIndex: "9999"
    };
  }

  function onClickOutside(e: MouseEvent) {
    if (!open.value) return;

    const target = e.target as Node;
    if (menuRef.value?.contains(target)) return;
    if (btnRef.value?.contains(target)) return;

    open.value = false;
  }

  // Close this menu when another menu is opened
  function onCloseAllMenus(e: Event) {
    const customEvent = e as CustomEvent;
    if (customEvent.detail !== props.id && open.value) {
      open.value = false;
    }
  }

  // Close menu on scroll
  function onScroll() {
    if (open.value) {
      open.value = false;
    }
  }

  onMounted(() => {
    document.addEventListener("click", onClickOutside);
    window.addEventListener("scroll", onScroll, true);
    window.addEventListener("resize", onScroll);
    document.addEventListener("keydown", onKeydown);
    window.addEventListener(CLOSE_ALL_MENUS_EVENT, onCloseAllMenus);
  });

  onBeforeUnmount(() => {
    document.removeEventListener("click", onClickOutside);
    window.removeEventListener("scroll", onScroll, true);
    window.removeEventListener("resize", onScroll);
    document.removeEventListener("keydown", onKeydown);
    window.removeEventListener(CLOSE_ALL_MENUS_EVENT, onCloseAllMenus);
  });

  function onKeydown(e: KeyboardEvent) {
    if (e.key === "Escape" && open.value) {
      open.value = false;
    }
  }

  function handleEditCompany() {
    emit("edit", props.id);
    open.value = false;
  }

  function handleEditContacts() {
    emit("edit-contacts", props.id);
    open.value = false;
  }

  function handleToggleStatus() {
    emit("toggle-status", props.id);
    open.value = false;
  }

  async function handleDelete() {
    open.value = false;
    
    const result = await Swal.fire({
      title: 'ยืนยันการลบ?',
      text: 'คุณต้องการลบบริษัทนี้ใช่หรือไม่? การดำเนินการนี้ไม่สามารถย้อนกลับได้',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#dc2626',
      cancelButtonColor: '#6b7280',
      confirmButtonText: 'ลบ',
      cancelButtonText: 'ยกเลิก'
    });
    
    if (result.isConfirmed) {
      emit("delete", props.id);
    }
  }
</script>

<template>
  <div class="inline-block">
    <button 
      class="btn btn-ghost btn-sm"
      @click.stop="toggleMenu"
      aria-haspopup="true"
      :aria-expanded="open"
      :aria-controls="'menu-'+props.id"
      ref="btnRef"
      :data-testid="`action-menu-btn-${props.id}`"
    >
      <Ellipsis class="w-4 h-4" />
    </button>

    <teleport to="body">
    <div
      v-if="open"
      ref="menuRef"
      :style="menuStyle"
      class="rounded-box shadow bg-base-100 w-44 menu p-2"
      data-theme="light"
      role="menu"
      :id="'menu-'+props.id"
    >
      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left rounded-lg text-sm"
        @click="handleEditCompany"
        role="menuitem"
        :data-testid="`action-edit-sponsor-${props.id}`"
      >
        <PenLine class="w-4 h-4" /> แก้ไขข้อมูลบริษัท
      </button>

      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left rounded-lg text-sm"
        @click="handleEditContacts"
        role="menuitem"
        :data-testid="`action-edit-contacts-${props.id}`"
      >
        <Users class="w-4 h-4" /> แก้ไขผู้ติดต่อ
      </button>

      <button
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left rounded-lg text-sm"
        @click="handleToggleStatus"
        role="menuitem"
        :data-testid="`action-toggle-status-${props.id}`"
      >
        <template v-if="props.status === 'active'">
          <Slash class="w-4 h-4"/> ปิดใช้งาน
        </template>
        <template v-else>
          <CheckCircle class="w-4 h-4" /> เปิดใช้งาน
        </template>
      </button>

      <div style="height:1px; background:rgba(0,0,0,0.06); margin:6px 0;"></div>

      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-red-50 text-red-500 w-full text-left rounded-lg text-sm"
        @click="handleDelete"
        role="menuitem"
        :data-testid="`action-delete-${props.id}`"
      >
        <Trash class="w-4 h-4" /> ลบบริษัท
      </button>
    </div>
    </teleport>
  </div>
</template>

