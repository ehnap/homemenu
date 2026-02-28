import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue'),
      meta: { public: true },
    },
    {
      path: '/',
      redirect: '/recipes',
    },
    {
      path: '/recipes',
      name: 'RecipeList',
      component: () => import('../views/RecipeList.vue'),
    },
    {
      path: '/recipes/new',
      name: 'RecipeCreate',
      component: () => import('../views/RecipeEdit.vue'),
    },
    {
      path: '/recipes/:id',
      name: 'RecipeDetail',
      component: () => import('../views/RecipeDetail.vue'),
    },
    {
      path: '/recipes/:id/edit',
      name: 'RecipeEdit',
      component: () => import('../views/RecipeEdit.vue'),
    },
    {
      path: '/meal-plans',
      name: 'MealPlanList',
      component: () => import('../views/MealPlanList.vue'),
    },
    {
      path: '/meal-plans/generate',
      name: 'MealPlanGenerate',
      component: () => import('../views/MealPlanGenerate.vue'),
    },
    {
      path: '/meal-plans/:id',
      name: 'MealPlanEditor',
      component: () => import('../views/MealPlanEditor.vue'),
    },
    {
      path: '/meal-plans/:id/shopping',
      name: 'ShoppingList',
      component: () => import('../views/ShoppingList.vue'),
    },
    {
      path: '/share/:token',
      name: 'ShareView',
      component: () => import('../views/ShareView.vue'),
      meta: { public: true },
    },
    {
      path: '/settings',
      name: 'Settings',
      component: () => import('../views/Settings.vue'),
    },
  ],
})

router.beforeEach((to) => {
  const token = localStorage.getItem('access_token')
  if (!to.meta.public && !token) {
    return { name: 'Login' }
  }
})

export default router
