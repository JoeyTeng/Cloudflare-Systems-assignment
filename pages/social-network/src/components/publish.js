import React, { useState } from "react";
import { useNavigate, Link } from "@reach/router";

const handleChange = (setter) => (e) => {
    setter(e.target.value);
};

const Post = ({ navigate }) => {
  const [title, setTitle] = useState("");
  const [username, setUsername] = useState("");
  const [content, setContent] = useState("");

  const handleSubmit = async () => {
    const resp = await fetch(
      `https://worker.joey-teng-dev.workers.dev/post`, {
      method: 'POST',
      cache: 'no-cache',
      headers: {
        'Content-Type': 'application/json',
      },
      redirect: 'follow',
      referrerPolicy: 'no-referrer',
      body: JSON.stringify({
        'title': title,
        'username': username,
        'content': content,
      })
    });
    if (resp.ok) {
      navigate(`/`);
    }
  };

  return (
    <div>
      <h1>Publish New Post</h1>
      <div className="mb-3">
        <label htmlFor="title" className="form-label">Title</label>
        <input type="text" className="form-control" id="title" onChange={handleChange(setTitle)} />
      </div>
      <div className="mb-3">
        <label htmlFor="username" className="form-label">Author Name</label>
        <input type="text" className="form-control" id="username" onChange={handleChange(setUsername)} />
      </div>
      <div className="mb-3">
        <label htmlFor="content" className="form-label">Body</label>
        <input type="text" className="form-control" id="content" onChange={handleChange(setContent)} />
      </div>
      <button type="submit" className="btn btn-primary login-button" onClick={handleSubmit}>Submit</button>
      <p>
        <Link to="/">Go back</Link>
      </p>
    </div>
  );
};

export default Post;
