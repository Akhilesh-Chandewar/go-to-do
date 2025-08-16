import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "../App";

export type Todo = {
    id: string;      // 👈 backend returns "id" (ObjectID)
    body: string;
    completed: boolean;
};

const TodoItem = ({ todo }: { todo: Todo }) => {
    const queryClient = useQueryClient();

    const { mutate: updateTodo, isPending: isUpdating } = useMutation({
        mutationKey: ["updateTodo", todo.id],
        mutationFn: async () => {
            if (todo.completed) return alert("Todo is already completed");

            const res = await fetch(BASE_URL + `/api/todos/${todo.id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ completed: true }),
            });
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Something went wrong");
            return data;
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
    });

    const { mutate: deleteTodo, isPending: isDeleting } = useMutation({
        mutationKey: ["deleteTodo", todo.id],
        mutationFn: async () => {
            const res = await fetch(BASE_URL + `/api/todos/${todo.id}`, {
                method: "DELETE",
            });
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || "Something went wrong");
            return data;
        },
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
    });

    return (
        <Flex gap={2} alignItems="center">
            <Flex
                flex={1}
                alignItems="center"
                border="1px"
                borderColor="gray.600"
                p={2}
                borderRadius="lg"
                justifyContent="space-between"
            >
                <Text
                    color={todo.completed ? "green.200" : "yellow.100"}
                    textDecoration={todo.completed ? "line-through" : "none"}
                >
                    {todo.body}
                </Text>
                {todo.completed ? (
                    <Badge ml="1" colorScheme="green">Done</Badge>
                ) : (
                    <Badge ml="1" colorScheme="yellow">In Progress</Badge>
                )}
            </Flex>

            <Flex gap={2} alignItems="center">
                <Box color="green.500" cursor="pointer" onClick={() => updateTodo()}>
                    {!isUpdating ? <FaCheckCircle size={20} /> : <Spinner size="sm" />}
                </Box>
                <Box color="red.500" cursor="pointer" onClick={() => deleteTodo()}>
                    {!isDeleting ? <MdDelete size={25} /> : <Spinner size="sm" />}
                </Box>
            </Flex>
        </Flex>
    );
};

export default TodoItem;
