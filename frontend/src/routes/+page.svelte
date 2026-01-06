<script lang="ts">
  import { Plus, Trash } from "@lucide/svelte";
  import { todoClient } from "$lib/api";
  import { type Todo } from "$lib/gen/todo/v1/todo_pb";
  import { onMount } from "svelte";

  let todos = $state<Todo[]>();
  let inputText = $state("");

  const loadTodos = async () => {
    const res = await todoClient.getTodos({});
    todos = res.todos;
    todos.sort((a, b) => {
      // Primary sort: unchecked (false) first
      if (a.completed !== b.completed) {
        return a.completed ? 1 : -1;
      }

      // Secondary sort: by date (newest first)
      return b.id - a.id;
    });
  };

  const createTodo = async () => {
    if (inputText.trim().length === 0) {
      return;
    }
    await todoClient.createTodo({
      text: inputText,
    });
    await loadTodos();
    inputText = "";
  };

  const updateTodo = async (id: number, completed: boolean) => {
    await todoClient.updateTodo({
      id,
      completed,
    });
    await loadTodos();
  };

  const deleteTodo = async (id: number) => {
    await todoClient.deleteTodo({
      id,
    });
    await loadTodos();
  };

  onMount(async () => {
    await loadTodos();
  });
</script>

<main class="max-w-2xl container mx-auto p-2">
  <h1 class="text-center font-bold text-2xl">todos</h1>
  <form class="flex flex-row gap-2">
    <input
      class="w-full input"
      type="text"
      bind:value={inputText}
      placeholder="type your todo here ..."
    />
    <button class="btn btn-primary" onclick={createTodo}
      ><Plus /> add todo</button
    >
  </form>
  <ul>
    {#each todos as todo}
      <li class="p-2">
        <input
          class="checkbox"
          type="checkbox"
          bind:checked={todo.completed}
          onchange={(e) => {
            updateTodo(todo.id, e.currentTarget.checked);
          }}
        />
        {todo.text}
        <button
          class="btn btn-accent btn-xs"
          onclick={() => {
            deleteTodo(todo.id);
          }}><Trash size="12" /></button
        >
      </li>
    {/each}
  </ul>
</main>
