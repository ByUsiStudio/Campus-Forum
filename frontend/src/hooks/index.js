import { ref, onMounted, onUnmounted } from 'vue'

export function usePagination(initialPage = 1, initialPageSize = 20) {
  const page = ref(initialPage)
  const pageSize = ref(initialPageSize)
  const totalPages = ref(1)
  const totalCount = ref(0)

  const setPage = (newPage) => {
    if (newPage >= 1 && newPage <= totalPages.value) {
      page.value = newPage
    }
  }

  const nextPage = () => {
    if (page.value < totalPages.value) {
      page.value++
    }
  }

  const prevPage = () => {
    if (page.value > 1) {
      page.value--
    }
  }

  const goToFirstPage = () => {
    page.value = 1
  }

  const goToLastPage = () => {
    page.value = totalPages.value
  }

  const resetPage = () => {
    page.value = 1
  }

  const setTotalPages = (total) => {
    totalPages.value = Math.max(1, total)
  }

  const setTotalCount = (count) => {
    totalCount.value = count
    totalPages.value = Math.max(1, Math.ceil(count / pageSize.value))
  }

  return {
    page,
    pageSize,
    totalPages,
    totalCount,
    setPage,
    nextPage,
    prevPage,
    goToFirstPage,
    goToLastPage,
    resetPage,
    setTotalPages,
    setTotalCount
  }
}

export function useLocalStorage(key, defaultValue) {
  const storedValue = localStorage.getItem(key)
  const data = ref(storedValue ? JSON.parse(storedValue) : defaultValue)

  const save = (value) => {
    data.value = value
    localStorage.setItem(key, JSON.stringify(value))
  }

  const clear = () => {
    data.value = defaultValue
    localStorage.removeItem(key)
  }

  return {
    data,
    save,
    clear
  }
}

export function useDebounce(value, delay = 300) {
  const debouncedValue = ref(value.value)

  let timeoutId = null

  const update = () => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => {
      debouncedValue.value = value.value
    }, delay)
  }

  onMounted(() => {
    update()
  })

  onUnmounted(() => {
    clearTimeout(timeoutId)
  })

  return debouncedValue
}

export function useThrottle(value, delay = 1000) {
  const throttledValue = ref(value.value)

  let lastTime = 0

  const update = () => {
    const now = Date.now()
    if (now - lastTime >= delay) {
      throttledValue.value = value.value
      lastTime = now
    }
  }

  onMounted(() => {
    update()
  })

  return throttledValue
}

export function useScroll() {
  const scrollX = ref(0)
  const scrollY = ref(0)
  const scrollTop = ref(0)
  const scrollHeight = ref(0)
  const clientHeight = ref(0)

  const handleScroll = () => {
    scrollX.value = window.scrollX || window.pageXOffset
    scrollY.value = window.scrollY || window.pageYOffset
    scrollTop.value = document.documentElement.scrollTop || document.body.scrollTop
    scrollHeight.value = document.documentElement.scrollHeight || document.body.scrollHeight
    clientHeight.value = window.innerHeight || document.documentElement.clientHeight
  }

  onMounted(() => {
    window.addEventListener('scroll', handleScroll, { passive: true })
    handleScroll()
  })

  onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
  })

  const isAtBottom = () => {
    return scrollTop.value + clientHeight.value >= scrollHeight.value - 100
  }

  const isAtTop = () => {
    return scrollTop.value < 50
  }

  const scrollToTop = (smooth = true) => {
    window.scrollTo({
      top: 0,
      behavior: smooth ? 'smooth' : 'auto'
    })
  }

  const scrollToElement = (selector, offset = 0) => {
    const element = document.querySelector(selector)
    if (element) {
      const elementTop = element.offsetTop
      window.scrollTo({
        top: elementTop + offset,
        behavior: 'smooth'
      })
    }
  }

  return {
    scrollX,
    scrollY,
    scrollTop,
    scrollHeight,
    clientHeight,
    isAtBottom,
    isAtTop,
    scrollToTop,
    scrollToElement
  }
}

export function useResize() {
  const width = ref(window.innerWidth)
  const height = ref(window.innerHeight)

  const handleResize = () => {
    width.value = window.innerWidth
    height.value = window.innerHeight
  }

  onMounted(() => {
    window.addEventListener('resize', handleResize, { passive: true })
  })

  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
  })

  const isMobile = () => width.value < 768
  const isTablet = () => width.value >= 768 && width.value < 1024
  const isDesktop = () => width.value >= 1024

  return {
    width,
    height,
    isMobile,
    isTablet,
    isDesktop
  }
}

export function useClickOutside(ref, callback) {
  const handleClick = (event) => {
    if (ref.value && !ref.value.contains(event.target)) {
      callback(event)
    }
  }

  onMounted(() => {
    document.addEventListener('click', handleClick)
  })

  onUnmounted(() => {
    document.removeEventListener('click', handleClick)
  })
}

export function useKeydown(key, callback) {
  const handleKeydown = (event) => {
    if (event.key === key) {
      callback(event)
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeydown)
  })
}

export function useIntersectionObserver(targetRef, options = {}) {
  const isIntersecting = ref(false)

  const observerCallback = (entries) => {
    entries.forEach(entry => {
      isIntersecting.value = entry.isIntersecting
    })
  }

  onMounted(() => {
    const observer = new IntersectionObserver(observerCallback, {
      root: options.root || null,
      rootMargin: options.rootMargin || '0px',
      threshold: options.threshold || 0
    })

    if (targetRef.value) {
      observer.observe(targetRef.value)
    }

    onUnmounted(() => {
      observer.disconnect()
    })
  })

  return isIntersecting
}

export function useFetch(url, options = {}) {
  const data = ref(null)
  const error = ref(null)
  const loading = ref(false)

  const fetchData = async () => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(url, options)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      data.value = await response.json()
    } catch (err) {
      error.value = err
    } finally {
      loading.value = false
    }
  }

  onMounted(() => {
    if (!options.lazy) {
      fetchData()
    }
  })

  return {
    data,
    error,
    loading,
    fetchData
  }
}

export function useInterval(callback, delay) {
  let intervalId = null

  const start = () => {
    if (intervalId) return
    intervalId = setInterval(callback, delay)
  }

  const stop = () => {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  const reset = () => {
    stop()
    start()
  }

  onMounted(() => {
    if (delay > 0) {
      start()
    }
  })

  onUnmounted(() => {
    stop()
  })

  return {
    start,
    stop,
    reset
  }
}

export function useTimeout(callback, delay) {
  let timeoutId = null

  const start = () => {
    if (timeoutId) return
    timeoutId = setTimeout(() => {
      callback()
      timeoutId = null
    }, delay)
  }

  const cancel = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  return {
    start,
    cancel
  }
}