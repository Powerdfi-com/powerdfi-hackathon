import { useMutation, useQuery } from "@tanstack/react-query";
import GhostContentAPI from "@tryghost/content-api";

export const ghost = new GhostContentAPI({
  url: "https://blog.powerdfi.com",
  key: "aa72e8254cdfe6a5f23f71e9de", // replace this with your API key
  version: "v4",
});

export const GhostAPI = {
  fetchBlogs: () =>
    useMutation({
      mutationKey: ["fetch blogs"],
      mutationFn: () => ghost.posts.browse({ limit: "all" }),
    }),
  fetchBlogBySlug: (slug: string) =>
    useQuery({
      queryKey: ["fetch blog by slug", slug],
      queryFn: () => ghost.posts.read({ slug }),
    }),
};

export const getPosts = async () => {
  try {
    return await ghost.posts.browse({
      limit: "all",
    });
  } catch (error) {
    console.error(error);
  }
};

export const getSinglePost = async (slug: string) => {
  try {
    return await ghost.posts.read({
      slug,
    });
  } catch (error) {
    console.error(error);
  }
};
