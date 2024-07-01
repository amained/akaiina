<script lang="ts">
	import getRequestClient from '$lib/getRequestClient';
	import auth from '../client';
	export let cookie;

	import { onMount } from 'svelte';
	import { Auth0Provider } from '$lib/auth';
	const params = new URLSearchParams();
	let auth0 = new Auth0Provider(null, cookie);
	if (!auth0.isAuthenticated()) {
		params.set('returnTo', window.location.pathname);
	}
</script>

<div>
	{#if auth0 != null}
		{#if auth0.isAuthenticated()}
			authenticated
		{:else if params.has('returnTo')}
			not authenticated
			<div class="authStatus">
				<form method="GET" action={'/login?' + params.toString()}>
					<button type="submit">sign in i guess</button>
				</form>
			</div>
		{:else}
			loading...
		{/if}
	{:else}
		loading
	{/if}
</div>
