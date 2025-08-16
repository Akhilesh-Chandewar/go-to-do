import { Flex, Spinner, Stack, Text } from "@chakra-ui/react";
import TodoItem, { type Todo } from "./TodoItem";
import { useQuery } from "@tanstack/react-query";
import { BASE_URL } from "../App";

const TodoList = () => {
    const { data, isLoading } = useQuery<Todo[]>({
        queryKey: ["todos"],
        queryFn: async () => {
            const res = await fetch(BASE_URL + "/api/todos");
            const json = await res.json();
            if (!res.ok) throw new Error(json.error || "Something went wrong");
            return json.data || [];
        },
    });

    return (
        <>
            <Text
                fontSize="4xl"
                textTransform="uppercase"
                fontWeight="bold"
                textAlign="center"
                my={2}
                bgGradient="linear(to-l, #0b85f8, #00ffff)"
                bgClip="text"
            >
                Today's Tasks
            </Text>

            {isLoading && (
                <Flex justifyContent="center" my={4}>
                    <Spinner size="xl" />
                </Flex>
            )}

            {!isLoading && data?.length === 0 && (
                <Stack alignItems="center" gap="3">
                    <Text fontSize="xl" textAlign="center" color="gray.500">
                        All tasks completed! 🎉
                    </Text>
                    <img src="/go.png" alt="Go logo" width={70} height={70} />
                </Stack>
            )}

            <Stack gap={3}>
                {data?.map((todo) => (
                    <TodoItem key={todo.id} todo={todo} />
                ))}
            </Stack>
        </>
    );
};

export default TodoList;
