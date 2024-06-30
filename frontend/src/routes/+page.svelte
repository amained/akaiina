<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import Client, { yume, Local } from '../client';
	const encoreClient = new Client(Local);
	let mutation = createMutation({
		mutationFn: (params: yume.NewDocumentParams) => encoreClient.yume.NewDocument(params),
		onSuccess: (data) => {
			console.log(data);
		},
		onError: (error) => {
			console.error(error);
		}
		// TODO: Invalidate shit
	});
	let userid: string,
		content: string,
		name: string = '';
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
<input bind:value={userid} placeholder="user id" />
<input bind:value={name} placeholder="document name" />
<input bind:value={content} placeholder="content" />

<button
	on:click={() => $mutation.mutate({ OWNER: userid, NAME: name, CONTENT_BASE64: btoa(content) })}
	>test</button
>
