// useFocusTrap.ts
import { ref, watch, onMounted, onBeforeUnmount, nextTick, type Ref } from 'vue';

type UseFocusTrapOptions = {
  onClose?: () => void;
  disableBackdropClose?: boolean;
};

const FOCUSABLE_SELECTOR = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled]):not([type="hidden"])',
  'textarea:not([disabled])',
  'select:not([disabled])',
  '[tabindex]:not([tabindex="-1"])'
].join(',');

export function useFocusTrap(isOpen: Ref<boolean>, options: UseFocusTrapOptions = {}) {
  const modalRef = ref<HTMLElement | null>(null);
  const dialogId = `dialog-${Math.random().toString(36).slice(2, 9)}`;
  let previouslyFocusedElement: HTMLElement | null = null;

  function isVisible(el: HTMLElement) {
    return !!(el.offsetParent || el.getClientRects().length);
  }

  function focusFirstElement(focusError = false) {
    const container = modalRef.value;
    if (!container) return;

    if (focusError) {
      const errorField = container.querySelector<HTMLElement>('[data-has-error="true"]');
      if (errorField && typeof errorField.focus === 'function') {
        errorField.focus();
        return;
      }
    }

    const el = container.querySelector<HTMLElement>(FOCUSABLE_SELECTOR);
    if (el && typeof el.focus === 'function') el.focus();
  }

  function onKeydown(e: KeyboardEvent) {
    if (!isOpen.value) return;

    if (e.key === 'Escape') {
      e.preventDefault();
      options.onClose?.();
      return;
    }

    if (e.key !== 'Tab') return;

    const container = modalRef.value;
    if (!container) return;

    const nodeList = container.querySelectorAll<HTMLElement>(FOCUSABLE_SELECTOR);
    const focusable = Array.from(nodeList).filter((el) => {
      if (el.hasAttribute('disabled')) return false;
      if (el.getAttribute('tabindex') === '-1') return false;
      return isVisible(el);
    });

    if (focusable.length === 0) return;

    const first = focusable[0] as HTMLElement;
    const last = focusable[focusable.length - 1] as HTMLElement;

    if (e.shiftKey) {
      if (document.activeElement === first) {
        e.preventDefault();
        last.focus();
      }
      return;
    }

    if (document.activeElement === last) {
      e.preventDefault();
      first.focus();
    }
  }

  function onBackdropClick(e: MouseEvent) {
    if (options.disableBackdropClose) return;
    if (e.target === e.currentTarget) {
      options.onClose?.();
    }
  }

  watch(
    isOpen,
    async (open) => {
      if (open) {
        previouslyFocusedElement = document.activeElement as HTMLElement | null;
        document.body.style.overflow = 'hidden';
        await nextTick();
        focusFirstElement();
      } else {
        document.body.style.overflow = '';
        if (previouslyFocusedElement?.focus) previouslyFocusedElement.focus();
        previouslyFocusedElement = null;
      }
    },
    { flush: 'post' }
  );

  onMounted(() => window.addEventListener('keydown', onKeydown));
  onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown));

  return {
    modalRef,
    dialogId,
    focusFirstElement,
    onBackdropClick,
  };
}
