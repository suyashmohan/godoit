<script lang="ts">
  import { Plus, Trash2 } from "@lucide/svelte";
  import { todoClient } from "$lib/api";
  import { type Todo } from "$lib/gen/todo/v1/todo_pb";
  import { onMount } from "svelte";
  import { fly, fade } from "svelte/transition";
  import { quintOut } from "svelte/easing";

  let todos = $state<Todo[]>([]);
  let inputText = $state("");

  const loadTodos = async () => {
    const res = await todoClient.getTodos({});
    todos = res.todos;
    todos.sort((a, b) => {
      if (a.completed !== b.completed) {
        return a.completed ? 1 : -1;
      }
      return b.id - a.id;
    });
  };

  const createTodo = async () => {
    if (inputText.trim().length === 0) return;
    await todoClient.createTodo({ text: inputText });
    await loadTodos();
    inputText = "";
  };

  const updateTodo = async (id: number, completed: boolean) => {
    await todoClient.updateTodo({ id, completed });
    await loadTodos();
  };

  const deleteTodo = async (id: number) => {
    await todoClient.deleteTodo({ id });
    await loadTodos();
  };

  onMount(async () => {
    await loadTodos();
  });
</script>

<main class="min-h-screen flex flex-col items-center py-12 px-4 gap-8">
  
  <!-- Header Section -->
  <header class="text-center space-y-2">
    <h1 class="text-6xl md:text-8xl font-black tracking-tighter text-neutral uppercase drop-shadow-[4px_4px_0px_rgba(255,142,158,1)]">
      GODOIT
    </h1>
    <p class="text-xl font-bold bg-secondary inline-block px-4 py-1 border-2 border-neutral shadow-[4px_4px_0px_0px_#1A1A1A] -rotate-2">
      Just get it done.
    </p>
  </header>

  <!-- Input Section -->
  <section class="w-full max-w-xl">
    <div class="bg-base-200 p-6 rounded-xl neo-box flex flex-col sm:flex-row gap-4 items-stretch">
      <div class="grow">
        <input
          class="neo-input text-lg font-bold placeholder:font-normal placeholder:text-gray-400"
          type="text"
          bind:value={inputText}
          placeholder="What needs doing?"
          onkeydown={(e) => e.key === 'Enter' && createTodo()}
        />
      </div>
      <button 
        class="btn btn-primary btn-lg neo-btn font-black h-auto min-h-0 py-3 sm:py-0"
        onclick={createTodo}
      >
        <Plus strokeWidth={3} /> ADD
      </button>
    </div>
  </section>

  <!-- Todo List -->
  <ul class="w-full max-w-xl flex flex-col gap-4">
    {#each todos as todo (todo.id)}
      <li 
        in:fly={{ y: 20, duration: 400, easing: quintOut }}
        out:fade={{ duration: 200 }}
        class="group flex items-center gap-4 p-4 bg-white rounded-lg border-2 border-neutral shadow-[3px_3px_0px_0px_#1A1A1A] transition-all hover:shadow-[5px_5px_0px_0px_#A3E4D7] hover:-translate-y-1"
        class:opacity-60={todo.completed}
      >
        <input
          class="neo-checkbox shrink-0"
          type="checkbox"
          bind:checked={todo.completed}
          onchange={(e) => updateTodo(todo.id, e.currentTarget.checked)}
        />
        
        <span 
          class="grow font-bold text-lg break-all transition-all duration-300"
          class:line-through={todo.completed}
          class:text-gray-400={todo.completed}
        >
          {todo.text}
        </span>

        <button
          class="btn btn-square btn-sm btn-ghost hover:bg-red-100 hover:text-red-600 border-2 border-transparent hover:border-neutral transition-all"
          onclick={() => deleteTodo(todo.id)}
          aria-label="Delete todo"
        >
          <Trash2 size="20" strokeWidth={2.5} />
        </button>
      </li>
    {/each}
  </ul>
  
</main>
