--
-- mariaQL database dump
--

\restrict IluJBFXSxNlwpsMVXYCbRIKuHUEh9rRvKNsW2n41ObO0a1qKoD0QpRmhAwbLxFg

-- Dumped from database version 16.10
-- Dumped by pg_dump version 16.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: invitationstatus; Type: TYPE; Schema: public; Owner: maria
--

CREATE TYPE public.invitationstatus AS ENUM (
    'pending',
    'accepted',
    'rejected',
    'expired'
);


ALTER TYPE public.invitationstatus OWNER TO maria;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cpm; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.cpm (
    taskid integer NOT NULL,
    parentprojectid integer,
    dependencies integer[],
    duration integer
);


ALTER TABLE public.cpm OWNER TO maria;

--
-- Name: cpmresult; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.cpmresult (
    projectid integer NOT NULL,
    result json
);


ALTER TABLE public.cpmresult OWNER TO maria;

--
-- Name: invitation; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.invitation (
    id integer NOT NULL,
    username character varying(255) NOT NULL,
    projectid integer NOT NULL,
    role character varying,
    status public.invitationstatus DEFAULT 'pending'::public.invitationstatus NOT NULL,
    hasread boolean DEFAULT false
);


ALTER TABLE public.invitation OWNER TO maria;

--
-- Name: invitation_id_seq; Type: SEQUENCE; Schema: public; Owner: maria
--

CREATE SEQUENCE public.invitation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.invitation_id_seq OWNER TO maria;

--
-- Name: invitation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maria
--

ALTER SEQUENCE public.invitation_id_seq OWNED BY public.invitation.id;


--
-- Name: managers; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.managers (
    managername character varying(255) NOT NULL,
    projectid integer NOT NULL
);


ALTER TABLE public.managers OWNER TO maria;

--
-- Name: pert; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.pert (
    parenttaskid integer NOT NULL,
    predecessortaskid integer,
    optimistic integer NOT NULL,
    pessimistic integer NOT NULL,
    mostlikely integer NOT NULL,
    parentprojectid integer
);


ALTER TABLE public.pert OWNER TO maria;

--
-- Name: pertresult; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.pertresult (
    projectid integer NOT NULL,
    result json
);


ALTER TABLE public.pertresult OWNER TO maria;

--
-- Name: projects; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.projects (
    projectid integer NOT NULL,
    projectname character varying(255) NOT NULL,
    projectdescription text,
    projectstartdate timestamp without time zone NOT NULL,
    projectduedate timestamp without time zone,
    ownername character varying(255) NOT NULL
);


ALTER TABLE public.projects OWNER TO maria;

--
-- Name: projects_projectid_seq; Type: SEQUENCE; Schema: public; Owner: maria
--

CREATE SEQUENCE public.projects_projectid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.projects_projectid_seq OWNER TO maria;

--
-- Name: projects_projectid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maria
--

ALTER SEQUENCE public.projects_projectid_seq OWNED BY public.projects.projectid;


--
-- Name: tasks; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.tasks (
    taskid integer NOT NULL,
    taskname character varying(255) NOT NULL,
    taskdescription text,
    taskstatus character varying(255) DEFAULT 'pending'::character varying NOT NULL,
    taskstartdate timestamp without time zone,
    taskduedate timestamp without time zone,
    parentprojectid integer NOT NULL,
    assignedusername character varying(255),
    approved boolean DEFAULT false,
    taskcompleteddate timestamp without time zone,
    taskapproveddate timestamp without time zone,
    createdby character varying(255),
    archieved boolean DEFAULT false
);


ALTER TABLE public.tasks OWNER TO maria;

--
-- Name: tasks_taskid_seq; Type: SEQUENCE; Schema: public; Owner: maria
--

CREATE SEQUENCE public.tasks_taskid_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tasks_taskid_seq OWNER TO maria;

--
-- Name: tasks_taskid_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maria
--

ALTER SEQUENCE public.tasks_taskid_seq OWNED BY public.tasks.taskid;


--
-- Name: updates; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.updates (
    id integer NOT NULL,
    projectid integer NOT NULL,
    msg text NOT NULL,
    createdat timestamp without time zone DEFAULT now(),
    createdby character varying(255) NOT NULL,
    targettype character varying(10) NOT NULL,
    targetusername character varying(255),
    hasread boolean DEFAULT false,
    CONSTRAINT updates_targettype_check CHECK (((targettype)::text = ANY ((ARRAY['all'::character varying, 'user'::character varying])::text[])))
);


ALTER TABLE public.updates OWNER TO maria;

--
-- Name: updates_id_seq; Type: SEQUENCE; Schema: public; Owner: maria
--

CREATE SEQUENCE public.updates_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.updates_id_seq OWNER TO maria;

--
-- Name: updates_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: maria
--

ALTER SEQUENCE public.updates_id_seq OWNED BY public.updates.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: maria
--

CREATE TABLE public.users (
    username character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    hashedpassword character(60) NOT NULL,
    created timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    userstatus boolean DEFAULT true
);


ALTER TABLE public.users OWNER TO maria;

--
-- Name: invitation id; Type: DEFAULT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.invitation ALTER COLUMN id SET DEFAULT nextval('public.invitation_id_seq'::regclass);


--
-- Name: projects projectid; Type: DEFAULT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.projects ALTER COLUMN projectid SET DEFAULT nextval('public.projects_projectid_seq'::regclass);


--
-- Name: tasks taskid; Type: DEFAULT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.tasks ALTER COLUMN taskid SET DEFAULT nextval('public.tasks_taskid_seq'::regclass);


--
-- Name: updates id; Type: DEFAULT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.updates ALTER COLUMN id SET DEFAULT nextval('public.updates_id_seq'::regclass);


--
-- Name: cpm cpm_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.cpm
    ADD CONSTRAINT cpm_pkey PRIMARY KEY (taskid);


--
-- Name: cpmresult cpmresult_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.cpmresult
    ADD CONSTRAINT cpmresult_pkey PRIMARY KEY (projectid);


--
-- Name: invitation invitation_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.invitation
    ADD CONSTRAINT invitation_pkey PRIMARY KEY (id);


--
-- Name: managers managers_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.managers
    ADD CONSTRAINT managers_pkey PRIMARY KEY (managername, projectid);


--
-- Name: pert pert_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pert
    ADD CONSTRAINT pert_pkey PRIMARY KEY (parenttaskid);


--
-- Name: pertresult pertresult_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pertresult
    ADD CONSTRAINT pertresult_pkey PRIMARY KEY (projectid);


--
-- Name: projects projects_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_pkey PRIMARY KEY (projectid);


--
-- Name: tasks tasks_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_pkey PRIMARY KEY (taskid);


--
-- Name: updates updates_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.updates
    ADD CONSTRAINT updates_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (username);


--
-- Name: cpm cpm_taskid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.cpm
    ADD CONSTRAINT cpm_taskid_fkey FOREIGN KEY (taskid) REFERENCES public.tasks(taskid) ON DELETE CASCADE;


--
-- Name: cpmresult cpmresult_projectid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.cpmresult
    ADD CONSTRAINT cpmresult_projectid_fkey FOREIGN KEY (projectid) REFERENCES public.projects(projectid);


--
-- Name: updates fk_creator; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.updates
    ADD CONSTRAINT fk_creator FOREIGN KEY (createdby) REFERENCES public.users(username);


--
-- Name: pert fk_parentprojectid; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pert
    ADD CONSTRAINT fk_parentprojectid FOREIGN KEY (parentprojectid) REFERENCES public.projects(projectid);


--
-- Name: cpm fk_parentprojectid; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.cpm
    ADD CONSTRAINT fk_parentprojectid FOREIGN KEY (parentprojectid) REFERENCES public.projects(projectid);


--
-- Name: updates fk_project; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.updates
    ADD CONSTRAINT fk_project FOREIGN KEY (projectid) REFERENCES public.projects(projectid) ON DELETE CASCADE;


--
-- Name: invitation fk_projectid; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.invitation
    ADD CONSTRAINT fk_projectid FOREIGN KEY (projectid) REFERENCES public.projects(projectid);


--
-- Name: updates fk_target_user; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.updates
    ADD CONSTRAINT fk_target_user FOREIGN KEY (targetusername) REFERENCES public.users(username);


--
-- Name: invitation fk_username; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.invitation
    ADD CONSTRAINT fk_username FOREIGN KEY (username) REFERENCES public.users(username);


--
-- Name: managers managers_managername_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.managers
    ADD CONSTRAINT managers_managername_fkey FOREIGN KEY (managername) REFERENCES public.users(username) ON DELETE CASCADE;


--
-- Name: managers managers_projectid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.managers
    ADD CONSTRAINT managers_projectid_fkey FOREIGN KEY (projectid) REFERENCES public.projects(projectid) ON DELETE CASCADE;


--
-- Name: pert pert_parenttaskid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pert
    ADD CONSTRAINT pert_parenttaskid_fkey FOREIGN KEY (parenttaskid) REFERENCES public.tasks(taskid) ON DELETE CASCADE;


--
-- Name: pert pert_predecessortaskid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pert
    ADD CONSTRAINT pert_predecessortaskid_fkey FOREIGN KEY (predecessortaskid) REFERENCES public.tasks(taskid) ON DELETE SET NULL;


--
-- Name: pertresult pertresult_projectid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.pertresult
    ADD CONSTRAINT pertresult_projectid_fkey FOREIGN KEY (projectid) REFERENCES public.projects(projectid);


--
-- Name: projects projects_ownername_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.projects
    ADD CONSTRAINT projects_ownername_fkey FOREIGN KEY (ownername) REFERENCES public.users(username) ON DELETE CASCADE;


--
-- Name: tasks tasks_assignedusername_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_assignedusername_fkey FOREIGN KEY (assignedusername) REFERENCES public.users(username) ON DELETE SET NULL;


--
-- Name: tasks tasks_createdby_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_createdby_fkey FOREIGN KEY (createdby) REFERENCES public.users(username) ON DELETE SET NULL;


--
-- Name: tasks tasks_parentprojectid_fkey; Type: FK CONSTRAINT; Schema: public; Owner: maria
--

ALTER TABLE ONLY public.tasks
    ADD CONSTRAINT tasks_parentprojectid_fkey FOREIGN KEY (parentprojectid) REFERENCES public.projects(projectid) ON DELETE CASCADE;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO maria;


--
-- Name: TABLE invitation; Type: ACL; Schema: public; Owner: maria
--

GRANT ALL ON TABLE public.invitation TO maria;


--
-- Name: SEQUENCE invitation_id_seq; Type: ACL; Schema: public; Owner: maria
--

GRANT ALL ON SEQUENCE public.invitation_id_seq TO maria;


--
-- Name: TABLE updates; Type: ACL; Schema: public; Owner: maria
--

GRANT ALL ON TABLE public.updates TO maria;


--
-- Name: SEQUENCE updates_id_seq; Type: ACL; Schema: public; Owner: maria
--

GRANT ALL ON SEQUENCE public.updates_id_seq TO maria;


--
-- Name: DEFAULT PRIVILEGES FOR SEQUENCES; Type: DEFAULT ACL; Schema: public; Owner: maria
--

ALTER DEFAULT PRIVILEGES FOR ROLE maria IN SCHEMA public GRANT ALL ON SEQUENCES TO maria;


--
-- Name: DEFAULT PRIVILEGES FOR FUNCTIONS; Type: DEFAULT ACL; Schema: public; Owner: maria
--

ALTER DEFAULT PRIVILEGES FOR ROLE maria IN SCHEMA public GRANT ALL ON FUNCTIONS TO maria;


--
-- Name: DEFAULT PRIVILEGES FOR TABLES; Type: DEFAULT ACL; Schema: public; Owner: maria
--

ALTER DEFAULT PRIVILEGES FOR ROLE maria IN SCHEMA public GRANT ALL ON TABLES TO maria;


--
-- mariaQL database dump complete
--

\unrestrict IluJBFXSxNlwpsMVXYCbRIKuHUEh9rRvKNsW2n41ObO0a1qKoD0QpRmhAwbLxFg

