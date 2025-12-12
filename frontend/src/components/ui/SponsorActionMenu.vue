<script setup lang="ts">
  import { ref, nextTick, onMounted, onBeforeUnmount } from "vue";
  import { CheckCircle, Ellipsis, PenLine, Slash, Trash, Users } from "lucide-vue-next";

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

  function toggleMenu(e: MouseEvent) {
    e.stopPropagation();
    open.value = !open.value;
    btnRef.value = e.currentTarget as HTMLElement;
    nextTick(positionMenu);
  }

  function positionMenu() {
    if (!btnRef.value || !menuRef.value) return;

    const btn = btnRef.value.getBoundingClientRect();
    const menu = menuRef.value.getBoundingClientRect();

    let left = btn.right - menu.width;
    let top = btn.bottom + 8;

    const pad = 8;
    if (left < pad) left = pad;
    if (left + menu.width > window.innerWidth - pad)
      left = window.innerWidth - menu.width - pad;
    if (top + menu.height > window.innerHeight - pad)
      top = btn.top - menu.height - 8;

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

  onMounted(() => {
    document.addEventListener("click", onClickOutside);
    window.addEventListener("scroll", positionMenu, true);
    window.addEventListener("resize", positionMenu);
    document.addEventListener("keydown", onKeydown);
  });

  onBeforeUnmount(() => {
    document.removeEventListener("click", onClickOutside);
    window.removeEventListener("scroll", positionMenu, true);
    window.removeEventListener("resize", positionMenu);
    document.removeEventListener("keydown", onKeydown);
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

  function handleDelete() {
    emit("delete", props.id);
    open.value = false;
  }
</script>

<template>
  <button 
    class="btn btn-ghost btn-sm"
    @click.stop="toggleMenu"
    aria-haspopup="true"
    :aria-expanded="open"
    :aria-controls="'menu-'+props.id"
    ref="btnRef"
  >
    <Ellipsis class="w-4 h-4" />
  </button>

  <teleport to="body">
    <div
      v-if="open"
      ref="menuRef"
      :style="menuStyle"
      class="rounded-box shadow bg-base-100 w-40 menu p-2"
      data-theme="light"
      role="menu"
      :id="'menu-'+props.id"
    >
      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left"
        @click="handleEditCompany"
        role="menuitem"
      >
        <PenLine class="w-4 h-4" /> Edit Sponsor
      </button>

      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left"
        @click="handleEditContacts"
        role="menuitem"
      >
        <Users class="w-4 h-4" /> Edit Contact
      </button>

      <button
        class="flex items-center gap-2 px-3 py-2 hover:bg-gray-100 w-full text-left"
        @click="handleToggleStatus"
        role="menuitem"
      >
        <template v-if="props.status === 'active'">
          <Slash class="w-4 h-4"/> Set Inactive
        </template>
        <template v-else>
          <CheckCircle class="w-4 h-4" /> Set Active
        </template>
      </button>

      <div style="height:1px; background:rgba(0,0,0,0.06); margin:6px 0;"></div>

      <button 
        class="flex items-center gap-2 px-3 py-2 hover:bg-red-50 text-red-500 w-full text-left"
        @click="handleDelete"
        role="menuitem"
      >
        <Trash class="w-4 h-4" /> Delete
      </button>
    </div>
  </teleport>
</template>
