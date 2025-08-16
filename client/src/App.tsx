import { Container } from "@chakra-ui/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import TodoForm from "./components/TodoForm";
import TodoList from "./components/TodoList";
import Navbar from "./components/Navbar";

export const BASE_URL = "http://localhost:8000";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <Navbar />
      <Container maxW="lg" py={10}>
        <TodoForm />
        <TodoList />
      </Container>
    </QueryClientProvider>
  );
}

export default App;
