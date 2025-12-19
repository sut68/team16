import { ref, computed, watch, type Ref } from 'vue'

interface UsePaginationOptions {
  defaultPage?: number
  defaultPerPage?: number
  maxButtons?: number
}

export function usePagination<T>(
  items: Ref<T[]>,
  options: UsePaginationOptions = {}
) {
  const {
    defaultPage = 1,
    defaultPerPage = 10,
    maxButtons = 7,
  } = options

  const page = ref(defaultPage)
  const perPage = ref(defaultPerPage)

  // Computed
  const total = computed(() => items.value.length)

  const totalPages = computed(() =>
    Math.max(1, Math.ceil(total.value / perPage.value))
  )

  const paged = computed(() => {
    const start = (page.value - 1) * perPage.value
    return items.value.slice(start, start + perPage.value)
  })

  const pages = computed(() => {
    const tp = totalPages.value
    if (tp <= maxButtons) {
      return Array.from({ length: tp }, (_, i) => i + 1)
    }

    const half = Math.floor(maxButtons / 2)
    let start = page.value - half
    let end = page.value + half

    if (start < 1) {
      start = 1
      end = maxButtons
    } else if (end > tp) {
      end = tp
      start = tp - maxButtons + 1
    }

    return Array.from({ length: end - start + 1 }, (_, i) => start + i)
  })

  // Watchers
  watch(totalPages, (tp) => {
    if (page.value > tp) page.value = tp
  })

  watch(perPage, () => {
    page.value = 1
  })

  // Actions
  function prevPage() {
    if (page.value > 1) page.value--
  }

  function nextPage() {
    if (page.value < totalPages.value) page.value++
  }

  function goToPage(p: number) {
    if (p >= 1 && p <= totalPages.value) {
      page.value = p
    }
  }

  function resetPage() {
    page.value = 1
  }

  return {
    page,
    perPage,
    total,
    totalPages,
    paged,
    pages,
    prevPage,
    nextPage,
    goToPage,
    resetPage,
  }
}