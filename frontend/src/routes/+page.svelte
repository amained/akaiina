<script lang="ts">
	import { createMutation } from '@tanstack/svelte-query';
	import type { CreateMutationResult } from '@tanstack/svelte-query';
	import type { PageData } from './$types';
	import Client, { yume, Local } from '../client';
	import LoginStatus from '../components/LoginStatus.svelte';
	import { onMount } from 'svelte';

	export let data: PageData;
	let encoreClient: Client | null = null;
	let loggedIn = false;
	let mutation: null | CreateMutationResult<yume.Document, Error, yume.NewDocumentParams, unknown>;
	let userid: string = '';
	let content: string = '';
	let name: string = '';

	onMount(() => {
		// try {
		if (data.page_server_data.auth_token !== null) {
			encoreClient = new Client(Local, { auth: data.page_server_data.auth_token });
			loggedIn = true;
		}
		// }
		// catch (e) {
		// 	loggedIn = false;
		// 	console.log('probably invalid token');
		// 	encoreClient = new Client(Local, {});
		// }
	});

	$: mutation =
		typeof loggedIn === 'boolean' && loggedIn
			? createMutation({
					mutationFn: (params: yume.NewDocumentParams) => encoreClient!.yume.NewDocument(params),
					onSuccess: (data) => {
						// TODO: Handle success, maybe update UI or show message
					},
					onError: (error) => {
						console.error('Mutation error:', error);
						// TODO: Handle error, show error message
					}
				})
			: null;

	function handleMutation() {
		if ($mutation != null) {
			$mutation.mutate({ OWNER: userid, NAME: name, CONTENT_BASE64: btoa(content) });
		}
	}
</script>

{#if mutation !== null}
	<h1>Welcome to SvelteKit</h1>
	<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
	<input bind:value={userid} placeholder="user id" />
	<input bind:value={name} placeholder="document name" />
	<input bind:value={content} placeholder="content" />

	<button on:click={handleMutation}>test</button>

	<LoginStatus cookie={data.page_server_data} />
{:else}
	loading...
{/if}
