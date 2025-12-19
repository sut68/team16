import { ref, computed, watch, type Ref } from 'vue'

interface UseSearchOptions {
  debounceMs?: number
}

export function useSearch<T>(
  items: Ref<T[]>,
  filterFn: (item: T, term: string) => boolean,
  options: UseSearchOptions = {}
) {
  const { debounceMs = 300 } = options

  const searchQuery = ref('')
  const q = ref('')

  let timer: ReturnType<typeof setTimeout> | null = null

  // Debounced search
  watch(searchQuery, (val) => {
    if (timer) clearTimeout(timer)
    timer = setTimeout(() => {
      q.value = val.trim().toLowerCase()
      timer = null
    }, debounceMs)
  })

  // Filtered items
  const filtered = computed(() => {
    const term = q.value
    if (!term) return items.value
    return items.value.filter(item => filterFn(item, term))
  })

  function clearSearch() {
    searchQuery.value = ''
    q.value = ''
  }

  return {
    searchQuery,
    q,
    filtered,
    clearSearch,
  }
}