import { createToast } from 'mosha-vue-toastify'
import NProgress from 'nprogress'
import { START_LOCATION, createRouter, createWebHistory } from 'vue-router'

import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/dashboard.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/auth/login.vue'),
    meta: { layout: 'auth' },
  },
  {
    path: '/users',
    name: 'users',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/users/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/users/create',
    name: 'users.create',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/users/create.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/users/:id',
    name: 'users.edit',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/users/edit.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/videos',
    name: 'videos',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/videos/index.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/videos/create',
    name: 'videos.create',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/videos/create.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/videos/:id',
    name: 'videos.edit',
    // route level code-splitting
    // this generates a separate chunk (About.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import('@/views/videos/edit.vue'),
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: '404',
    component: () => import('@/views/errors/404.vue'),
  },
]
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: routes,
})
router.beforeEach(async (to, from) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.getters.isLoggedIn) {
    createToast(
      {
        title: 'Auth',
        description: 'Sorry, you should log in to access this section',
      },
      {
        type: 'danger',
        showIcon: true,
        hideProgressBar: true,
      },
    )

    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }
  if (to.meta.requiresGuest && userStore.getters.isLoggedIn) {
    return {
      name: 'home',
    }
  }
})
router.beforeResolve((to, from, next) => {
  if (to.name) {
    NProgress.start()
  }
  next()
})

router.afterEach(() => {
  NProgress.done()
})
export default router
