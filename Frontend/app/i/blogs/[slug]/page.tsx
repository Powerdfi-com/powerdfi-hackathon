import { Metadata, ResolvingMetadata } from "next";
import GhostBlog from "./blog";

export async function generateMetadata({ params }: { params: { slug: string } }, parent: ResolvingMetadata): Promise<Metadata> {
  const data = (await (await fetch(`https://blog.powerdfi.com/ghost/api/content/posts/slug/${params.slug}/?key=aa72e8254cdfe6a5f23f71e9de`)).json()).posts[0];
  return {
    title: data.meta_title,
    description: data.meta_description,
    openGraph: {
      title: data.og_title || data.title,
      description: data.og_description || data.excerpt,
      images: data.og_image || data.feature_image!
    }
  }
}

export default function Blog({ params }: { params: { slug: string } }) {
  return <GhostBlog slug={params.slug} />
}
