import { test } from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as Vue from 'vue'
import { safeProductDraft } from '../../resources/js/design/productDraft.ts'
const source = readFileSync('resources/js/Pages/Dashboard/Index.vue', 'utf8')
const meta = (q = '', page = 1, role = '') => ({ page, per_page: 25, total: 101, filtered_total: 76, last_page: 4, from: (page-1)*25+1, to: Math.min(page*25,76), filters: { q, role, verification:'', sort: 'newest' } })
function setup(extra = {}) {
    const gets = [], posts = [], timers = [], fetches = [], cleanups = []
    const props = Vue.reactive({ currentUser: { id: 'm', role: 'merchant' }, stats: {}, products: [], users: [{ id: 'o', name: 'SOKO-1', role: 'buyer', total_amount: 20, items: [{ id: 'line', quantity: 1, unit_price: 20 }] }], userPage: meta(), userVisibility: 'supplier', recordCapabilities: { users: { o: { cancel: false } } }, capabilities: { manageUsers: true, readProducts: true }, ...extra })
    const page = Vue.reactive({ url: '/dashboard?section=users', props: {} })
    let script = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    script += '\n globalThis.h = { userFilters, userFiltersPending, userListUrl, applyUserFilters, resetUserFilters, changeUserPage, changeUserPageSize, sortUserColumn, userListLoading, userListError, visibleRecords, userColumns, openRecord, detailRecord, detailKind, pendingOperation, refreshSelectedUser, accessSettingsHref, operationActions, confirmOperation, formError, busy };'
    const context = { exports: {}, defineProps: () => props, URL, URLSearchParams, AbortController,
        setTimeout: (callback, delay) => { timers.push({ callback, delay, cleared: false }); return timers.length }, clearTimeout: id => { if(timers[id-1]) timers[id-1].cleared = true },
        fetch: (url, options) => new Promise(resolve => fetches.push({ url, options, resolve })),
        require: name => name === 'vue' ? { ...Vue, onBeforeUnmount: fn => cleanups.push(fn) } : name === '@inertiajs/vue3' ? { usePage: () => page, router: { get: (...args) => gets.push(args), post: (...args) => posts.push(args), restore: () => undefined, remember: () => {} } } : name.includes('productDraft') ? { safeProductDraft } : name.includes('useCurrency') ? { useCurrency: () => ({ formatPrice: String }) } : {},
    }
    vm.runInNewContext(ts.transpileModule(script, { compilerOptions: { target: ts.ScriptTarget.ES2022, module: ts.ModuleKind.CommonJS } }).outputText, context)
    return { ...context.h, props, page, gets, posts, timers, fetches, cleanup: () => cleanups.forEach(fn => fn()) }
}
const flush = async () => { await Vue.nextTick(); await new Promise(resolve => setImmediate(resolve)) }
const respond = (request, record, cancel = false, visibility = 'supplier') => request.resolve({ ok: true, json: async () => ({ props: { selectedUser: record, selectedUserStatus: record ? 'available' : 'unavailable', userVisibility: visibility, recordCapabilities: { users: { o: { cancel } } } } }) })

test('Accounts use authoritative scoped rows and encoded role/verification filters',()=>{
 const h=setup({userPage:meta('ÉCOLE',3,'buyer')})
 h.userFilters.verification='unverified';h.applyUserFilters()
 const q=new URL(h.gets[0][0],'https://example.com').searchParams
 assert.equal(q.get('section'),'users');assert.equal(q.get('role'),'buyer');assert.equal(q.get('verification'),'unverified');assert.equal(q.get('q'),'ÉCOLE');assert.equal(q.get('page'),'1')
 assert.equal(h.visibleRecords.value.length,1);assert.equal(h.operationActions('users',h.props.users[0]).length,1);assert.equal(h.accessSettingsHref.value,'');h.cleanup()
})
test('Account search is serialized and immediate clear preserves newer input',async()=>{
 const h=setup();h.userFilters.q='first';await flush();assert.equal(h.timers.at(-1).delay,350);h.timers.at(-1).callback()
 h.userFilters.q='';await flush();assert.equal(h.timers.at(-1).delay,0)
 h.props.userPage=meta('first');assert.equal(h.userFilters.q,'');h.gets[0][2].onFinish();await flush()
 assert.equal(h.gets.length,2);assert.equal(new URL(h.gets[1][0],'https://example.com').searchParams.has('q'),false);h.cleanup()
})
test('Account sorting and page size reset pagination; pending filters block page navigation',()=>{
 const h=setup({userPage:meta('buyer',3)})
 h.userFilters.role='merchant';h.changeUserPage(4);assert.equal(h.gets.length,0)
 h.userFilters.role='';h.sortUserColumn({key:'name',direction:'asc'});assert.equal(h.gets.length,0)
 h.sortUserColumn({key:'created_at',direction:'asc'});assert.match(h.gets[0][0],/sort=oldest/);assert.match(h.gets[0][0],/page=1/)
 h.gets[0][2].onFinish();h.changeUserPageSize(50);assert.match(h.gets[1][0],/per_page=50/);h.cleanup()
})
test('Account selected lookup refreshes off-page identity without exposing edit actions',async()=>{
 const h=setup();h.openRecord('users',h.props.users[0]);const loading=h.refreshSelectedUser()
 assert.match(h.fetches[0].url,/selected_user=o/)
 respond(h.fetches[0],{id:'o',name:'Updated account',role:'buyer'});await loading
 assert.equal(h.detailRecord.value.name,'Updated account');assert.equal(h.operationActions('users',h.detailRecord.value).length,1);assert.equal(h.posts.length,0);h.cleanup()
})
test('Account denied or missing selected refresh reports error without replacing displayed identity',async()=>{
 for(const unavailable of [true,false]){
 const h=setup();h.openRecord('users',h.props.users[0]);const loading=h.refreshSelectedUser()
 if(unavailable)respond(h.fetches[0],null);else h.fetches[0].resolve({ok:false,status:403})
 await loading;assert.equal(h.detailRecord.value.id,'o');assert.ok(h.formError.value);assert.equal(h.posts.length,0);h.cleanup()
 }
})
test('Leaving Accounts aborts selected lookup and removes old details',async()=>{
 const h=setup();h.openRecord('users',h.props.users[0]);const loading=h.refreshSelectedUser()
 h.page.url='/dashboard?section=products';await flush();assert.equal(h.fetches[0].options.signal.aborted,true);assert.equal(h.detailRecord.value,null)
 respond(h.fetches[0],{id:'o',name:'Stale'});await loading;assert.equal(h.detailRecord.value,null);assert.equal(h.busy.value,false);h.cleanup()
})
test('Account access link requires the explicit authorized workspace navigation entry',()=>{
 const h=setup();h.page.props.workspaceNavigation=[{href:'/dashboard/settings?panel=payments'}];assert.equal(h.accessSettingsHref.value,'')
 h.page.props.workspaceNavigation=[{href:'/dashboard/settings?panel=access'}];assert.equal(h.accessSettingsHref.value,'/dashboard/settings?panel=access');h.cleanup()
})
