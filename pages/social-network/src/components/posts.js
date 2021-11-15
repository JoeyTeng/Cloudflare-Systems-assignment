import React, { useEffect, useState } from "react";
import { Link } from "@reach/router";

const Posts = () => {
  const [posts, setPosts] = useState([]);

  useEffect(() => {
    const getPosts = async () => {
      const resp = await fetch(
        "https://worker.joey-teng-dev.workers.dev/post",
      );
      console.log(resp);
      const postsResp = await resp.json();
      // as data from server is in ascending order
      setPosts(postsResp.map((post) => JSON.parse(post)).reverse());
    };

    getPosts();
  }, []);

  return (
    <div>
      <h1>Posts</h1>
      <Link to="/publish">Publish your post</Link>
      {posts.map((post) => (
        <div>
          <p>
            <h3>{post.title}</h3>
            user: {post.username} <br></br>
            body: {post.content}
          </p>
          <hr></hr>
        </div>
      ))}
    </div>
  );
};

export default Posts;
